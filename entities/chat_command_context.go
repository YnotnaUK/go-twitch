package entities

type ChatCommandContext struct {
	CommandName   string
	CommandParams []string
	Message       *ChatPrivateMessage
	Say           func(channel string, message string)
}
