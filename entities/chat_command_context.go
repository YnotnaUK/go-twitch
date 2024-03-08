package entities

type ChatCommandContext struct {
	CommandName   string
	CommandParams []string
	Message       *ChatPrivateMessage
	Reply         func(message *ChatPrivateMessage, response string)
	Say           func(channel string, message string)
}
