package entities

type ChatCommandContext struct {
	CommandName   string
	CommandParams []string
	Message       *ChatPrivateMessage
}
