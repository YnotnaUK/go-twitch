package entities

type ChatPrivateMessage struct {
	Channel  string
	Message  string
	Tags     map[string]string
	Username string
}
