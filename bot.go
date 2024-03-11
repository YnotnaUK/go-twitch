package twitch

import (
	"log"
	"strings"
)

type TwitchBot struct {
	chat              *ChatClient
	chatCommandPrefix string
	chatCommands      map[string][]ChatCommander
}

func (b *TwitchBot) ChatJoin(channel string) error {
	err := b.chat.Join(channel)
	if err != nil {
		return err
	}
	return nil
}

func (b *TwitchBot) ChatReply(message *ChatPrivateMessage, response string) {
	b.chat.Reply(message, response)
}

func (b *TwitchBot) ChatSay(channel string, message string) {
	b.chat.Say(channel, message)
}

func (b *TwitchBot) OnChatCommand(commandName string, command ChatCommander) {
	b.chatCommands[commandName] = append(b.chatCommands[commandName], command)
}

func (b *TwitchBot) OnChatConnect(handler func(message *ChatConnectMessage)) error {
	b.chat.OnConnect(handler)
	return nil
}

func (b *TwitchBot) OnChatJoin(handler func(message *ChatJoinMessage)) {
	b.chat.OnJoin(handler)
}

func (b *TwitchBot) OnChatPart(handler func(message *ChatPartMessage)) {
	b.chat.OnPart(handler)
}

func (b *TwitchBot) OnChatPing(handler func(message *ChatPingMessage)) {
	b.chat.OnPing(handler)
}

func (b *TwitchBot) OnChatPong(handler func(message *ChatPongMessage)) {
	b.chat.OnPong(handler)
}

func (b *TwitchBot) OnChatPrivateMessage(handler func(message *ChatPrivateMessage)) error {
	b.chat.OnPrivateMessage(handler)
	return nil
}

func (b *TwitchBot) Start() {
	// Create command handler
	b.chat.OnPrivateMessage(func(message *ChatPrivateMessage) {
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
					commandContext := &ChatCommandContext{}
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
	log.Println("Starting bot...")
	b.chat.Start()
}

func NewBot(authProvider AuthProvider) (*TwitchBot, error) {
	// Create chat client
	chat, err := NewChatClient(authProvider)
	if err != nil {
		return nil, err
	}
	// Create bot
	bot := &TwitchBot{
		chat:              chat,
		chatCommands:      make(map[string][]ChatCommander),
		chatCommandPrefix: "!",
	}
	return bot, nil
}
