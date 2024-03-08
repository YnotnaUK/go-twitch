package bot

import (
	"log"
	"strings"

	"github.com/ynotnauk/go-twitch/chat"
	"github.com/ynotnauk/go-twitch/entities"
	"github.com/ynotnauk/go-twitch/interfaces"
)

type Bot struct {
	chat              *chat.Client
	chatCommandPrefix string
	chatCommands      map[string][]interfaces.ChatCommander
}

func (b *Bot) ChatJoin(channel string) error {
	err := b.chat.Join(channel)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) ChatReply(message *entities.ChatPrivateMessage, response string) {
	b.chat.Reply(message, response)
}

func (b *Bot) ChatSay(channel string, message string) {
	b.chat.Say(channel, message)
}

func (b *Bot) OnChatCommand(commandName string, command interfaces.ChatCommander) {
	b.chatCommands[commandName] = append(b.chatCommands[commandName], command)
}

func (b *Bot) OnChatJoin(handler func(message *entities.ChatJoinMessage)) {
	b.chat.OnJoin(handler)
}

func (b *Bot) OnChatPrivateMessage(handler func(message *entities.ChatPrivateMessage)) error {
	b.chat.OnPrivateMessage(func(message *entities.ChatPrivateMessage) {
		// Check to see if a command has requested
		if strings.HasPrefix(message.Message, b.chatCommandPrefix) && len(message.Message) > 1 {
			messageParts := strings.Split(message.Message, " ")
			commandName := strings.TrimPrefix(messageParts[0], b.chatCommandPrefix)
			// Check if handler(s) have been loaded for the command
			handlers, ok := b.chatCommands[commandName]
			if ok {
				commandParams := messageParts[1:]
				// Ensure there is at least 1 command handler
				if len(handlers) > 0 {
					// Build command context
					commandContext := &entities.ChatCommandContext{}
					commandContext.CommandName = commandName
					if len(messageParts) > 1 {
						commandContext.CommandParams = commandParams
					}
					commandContext.Message = message
					commandContext.Reply = b.ChatReply
					commandContext.Say = b.ChatSay
					// Call each handler
					for _, handler := range handlers {
						handler.Execute(commandContext)
					}
				}
			}
		}
	})
	b.chat.OnPrivateMessage(handler)
	return nil
}

func (b *Bot) OnTwitchChatConnect(handler func(message *entities.ChatConnectMessage)) error {
	b.chat.OnConnect(handler)
	return nil
}

func (b *Bot) Start() {
	log.Println("Starting bot...")
	b.chat.Start()
}

func New(authProvider interfaces.AuthProvider) (*Bot, error) {
	// Create chat client
	chat, err := chat.NewClient(authProvider)
	if err != nil {
		return nil, err
	}
	// Create bot
	bot := &Bot{
		chat:              chat,
		chatCommands:      make(map[string][]interfaces.ChatCommander),
		chatCommandPrefix: "!",
	}
	return bot, nil
}
