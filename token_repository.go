package twitch

type TokenRepository interface {
	GetTokens() (*Tokens, error)
	SaveTokens(tokens *Tokens) error
}
