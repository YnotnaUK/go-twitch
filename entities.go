package twitch

type AuthRecord struct {
	AccessToken  string   `json:"accessToken"`
	ClientId     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	ExpiresIn    int64    `json:"expiresIn"`
	Login        string   `json:"login"`
	RefreshToken string   `json:"refreshToken"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"tokenType"`
	UserId       string   `json:"userId"`
}

type ChatCommandContext struct {
	CommandName   string
	CommandParams []string
	Message       *ChatPrivateMessage
	Reply         func(message *ChatPrivateMessage, response string)
	Say           func(channel string, message string)
}

type ChatConnectMessage struct {
	Hostname string
}

type ChatJoinMessage struct {
	Channel  string
	Username string
}

type ChatPartMessage struct {
	Channel  string
	Username string
}

type ChatPingMessage struct{}

type ChatPongMessage struct {
	Server    string
	Timestamp int64
}

type ChatPrivateMessage struct {
	Channel  string
	Message  string
	Tags     map[string]string
	Username string
}

type IrcMessage struct {
	Command string
	Raw     string
	Params  []string
	Source  *IrcMessageSource
	Tags    map[string]string
}

type IrcMessageSource struct {
	Nickname string
	Username string
	Host     string
}

type RefreshTokenFailed struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type RefreshTokenSuccess struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int64    `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type ValidateTokenFailed struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ValidateTokenSuccess struct {
	ClientId  string   `json:"client_id"`
	ExpiresIn int64    `json:"expires_in"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserId    string   `json:"user_id"`
}
