package chat

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ynotnauk/go-twitch/entities"
	"github.com/ynotnauk/go-twitch/interfaces"
)

const (
	serverAddress string = "irc.chat.twitch.tv:6697"
)

var (
	ErrBlankChannel       error = errors.New("channel cannot be blank")
	ErrBlankRawIrcMessage error = errors.New("rawIrcMessage cannot be blank")
)

type Client struct {
	authProvider               interfaces.AuthProvider
	connectionOutgoingChannel  chan string
	connectionIncommingChannel chan string
	disconnectChannel          chan bool
	onConnect                  func(message *entities.ChatConnectMessage)
	onPing                     func(message *entities.ChatPingMessage)
}

func (c *Client) connect() error {
	log.Printf("Attempting to connect to Twitch [%s]", serverAddress)

	// Create a dialer
	netDialer := &net.Dialer{
		KeepAlive: time.Second * 10,
	}

	// tls configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Attempt to connect to server
	connection, err := tls.DialWithDialer(netDialer, "tcp", serverAddress, tlsConfig)
	if err != nil {
		return err
	}
	log.Println("Connected to Twitch!")

	// Start all required go routines
	wg := &sync.WaitGroup{}
	wg.Add(3)
	c.startMessageParser(wg)
	c.startConnectionReader(wg, connection)
	c.startConnectionWriter(wg, connection)

	// Get login details from auth provider
	login, accessToken, err := c.authProvider.GetLoginAndAccessToken()
	if err != nil {
		return err
	}

	// Setup the connection
	c.send("CAP REQ :twitch.tv/commands twitch.tv/membership twitch.tv/tags")
	c.send("PASS oauth:" + accessToken)
	c.send("NICK " + login)

	// Wait for all go routines to close
	wg.Wait()
	log.Println("Disconnected from Twitch")
	return err
}

func (c *Client) handleParsedIrcMessage(parsedIrcMessage *entities.IrcMessage) error {
	switch parsedIrcMessage.Command {
	case "001":
		connectMessage := &entities.ChatConnectMessage{
			Hostname: serverAddress,
		}
		if c.onConnect != nil {
			c.onConnect(connectMessage)
		}
	case "PING":
		response := fmt.Sprintf("PONG :%s", parsedIrcMessage.Params[0])
		c.send(response)
		pingMessage := &entities.ChatPingMessage{}
		if c.onPing != nil {
			c.onPing(pingMessage)
		}
	default:
		log.Print("Unhandled command!")
		log.Print("================================================================================")
		log.Printf("%+v", parsedIrcMessage)
		return nil
	}
	// TODO: this is for testing so i can see what commands were handled
	log.Printf("Command Handled: %s", parsedIrcMessage.Command)
	return nil
}

func (c *Client) Join(channel string) error {
	// TODO: check if channel can be split by commas
	// If it can we need to check each segment
	if channel == "" {
		return ErrBlankChannel
	}
	// If channel does not start with a # add it
	if !strings.HasPrefix(channel, "#") {
		channel = fmt.Sprintf("#%s", channel)
	}
	c.send(fmt.Sprintf("JOIN %s", channel))
	return nil
}

func (c *Client) OnConnect(handler func(message *entities.ChatConnectMessage)) {
	c.onConnect = handler
}

func (c *Client) OnPing(handler func(message *entities.ChatPingMessage)) {
	c.onPing = handler
}

func (c *Client) parseRawIrcMessage(rawIrcMessage string) (*entities.IrcMessage, error) {
	// Ensure rawIrcMessage is not blank
	if rawIrcMessage == "" {
		return nil, ErrBlankRawIrcMessage
	}
	// Create parsed IrcMessage struct
	parsedIrcMessage := &entities.IrcMessage{
		Raw: rawIrcMessage,
	}
	// Split the raw irc message into sections
	rawIrcMessageSplit := strings.Split(rawIrcMessage, " ")
	rawIrcMessageSplitIndex := 0
	// Check if the first section is a tags section
	if strings.HasPrefix(rawIrcMessageSplit[rawIrcMessageSplitIndex], "@") {
		parsedIrcMessageTags, err := c.parseRawIrcMessageTags(rawIrcMessageSplit[rawIrcMessageSplitIndex])
		if err != nil {
			return nil, err
		}
		parsedIrcMessage.Tags = parsedIrcMessageTags
		rawIrcMessageSplitIndex++
	}
	// Message source
	if strings.HasPrefix(rawIrcMessageSplit[rawIrcMessageSplitIndex], ":") {
		parsedIrcMessageSource, err := c.parseRawIrcMessageSource(rawIrcMessageSplit[rawIrcMessageSplitIndex])
		if err != nil {
			return nil, err
		}
		parsedIrcMessage.Source = parsedIrcMessageSource
		rawIrcMessageSplitIndex++
	}
	// Message command
	parsedIrcMessage.Command = rawIrcMessageSplit[rawIrcMessageSplitIndex]
	rawIrcMessageSplitIndex++
	// Remaining segments added to params
	var parsedParams []string
	for paramIndex, paramValue := range rawIrcMessageSplit[rawIrcMessageSplitIndex:] {
		// If its the first index we want to remove the : if it exists
		if paramIndex == 0 {
			paramValue = strings.TrimPrefix(paramValue, ":")
		}
		parsedParams = append(parsedParams, paramValue)
	}
	parsedIrcMessage.Params = parsedParams
	return parsedIrcMessage, nil
}

func (c *Client) parseRawIrcMessageSource(rawIrcMessageSource string) (*entities.IrcMessageSource, error) {
	parsedIrcMessageSource := &entities.IrcMessageSource{}
	rawIrcMessageSource = strings.TrimPrefix(rawIrcMessageSource, ":")
	regex := regexp.MustCompile(`!|@`)
	rawIrcMessageSourceSplit := regex.Split(rawIrcMessageSource, -1)
	switch len(rawIrcMessageSourceSplit) {
	case 1:
		parsedIrcMessageSource.Host = rawIrcMessageSourceSplit[0]
	case 2:
		parsedIrcMessageSource.Nickname = rawIrcMessageSourceSplit[0]
		parsedIrcMessageSource.Host = rawIrcMessageSourceSplit[1]
	default:
		parsedIrcMessageSource.Nickname = rawIrcMessageSourceSplit[0]
		parsedIrcMessageSource.Username = rawIrcMessageSourceSplit[1]
		parsedIrcMessageSource.Host = rawIrcMessageSourceSplit[2]
	}
	return parsedIrcMessageSource, nil
}

func (c *Client) parseRawIrcMessageTags(rawIrcMessageTags string) (map[string]string, error) {
	parsedIrcMessageTags := make(map[string]string)
	rawIrcMessageTags = strings.TrimPrefix(rawIrcMessageTags, "@")
	for _, rawIrcMessageTag := range strings.Split(rawIrcMessageTags, ";") {
		rawIrcMessageTagPairs := strings.SplitN(rawIrcMessageTag, "=", 2)
		rawIrcMessageTagKey := rawIrcMessageTagPairs[0]
		rawIrcMessageTagValue := rawIrcMessageTagPairs[1]
		parsedIrcMessageTags[rawIrcMessageTagKey] = rawIrcMessageTagValue
	}
	return parsedIrcMessageTags, nil
}

func (c *Client) send(line string) {
	// TODO: below 2 lines are for testing and need to removed at some point
	log.Println("Sending: " + line)
	line = line + "\r\n"
	c.connectionOutgoingChannel <- line
}

func (c *Client) Start() error {
	log.Println("Starting chat client")
	for {
		err := c.connect()
		switch err {
		default:
			return err
		}
	}
}

func (c *Client) startConnectionReader(wg *sync.WaitGroup, connection io.Reader) {
	log.Println("Starting connection reader")
	go func() {
		defer func() {
			close(c.disconnectChannel)
			log.Println("Connection reader has closed")
			wg.Done()
		}()
		tp := textproto.NewReader(bufio.NewReader(connection))
		for {
			// Check if there is a new line to read
			line, err := tp.ReadLine()
			if err != nil {
				return
			}
			// Split line to make sure no multiple messages per line
			rawIrcMessages := strings.Split(line, "\r\n")
			// Loop over and pass each message to the connectionIncommingChannel
			for _, rawIrcMessage := range rawIrcMessages {
				c.connectionIncommingChannel <- rawIrcMessage
			}
		}
	}()
}

func (c *Client) startConnectionWriter(wg *sync.WaitGroup, connection io.Writer) {
	log.Println("Starting connection writer")
	go func() {
		defer func() {
			log.Println("Connection writer has closed")
			wg.Done()
		}()
		for {
			select {
			case <-c.disconnectChannel:
				return
			case rawIrcMessage := <-c.connectionOutgoingChannel:
				connection.Write([]byte(rawIrcMessage))
			}
		}
	}()
}

func (c *Client) startMessageParser(wg *sync.WaitGroup) {
	log.Println("Starting message parser")
	go func() {
		defer func() {
			log.Println("message parser has closed")
			wg.Done()
		}()
		for {
			select {
			case <-c.disconnectChannel:
				return
			case rawIrcMessage := <-c.connectionIncommingChannel:
				parsedIrcMessage, err := c.parseRawIrcMessage(rawIrcMessage)
				if err != nil {
					log.Panic(err)
				}
				err = c.handleParsedIrcMessage(parsedIrcMessage)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}()
}

func NewClient(authProvider interfaces.AuthProvider) (*Client, error) {
	client := &Client{
		authProvider:               authProvider,
		connectionOutgoingChannel:  make(chan string, 64),
		connectionIncommingChannel: make(chan string, 64),
		disconnectChannel:          make(chan bool),
	}
	return client, nil
}
