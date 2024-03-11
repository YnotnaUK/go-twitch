package twitch

type AuthProvider interface {
	GetAccessToken() (string, error)
	GetLoginAndAccessToken() (string, string, error)
}

type AuthStorer interface {
	GetByUserId(userId string) (*AuthRecord, error)
	UpdateByUserId(auth *AuthRecord) error
}

type ChatCommander interface {
	Execute(message *ChatCommandContext)
}
