package entities

type IrcMessage struct {
	Command string
	Raw     string
	Params  []string
	Source  *IrcMessageSource
	Tags    map[string]string
}
