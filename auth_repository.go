package twitch

type AuthRepository interface {
	Load() (*Auth, error)
	Save(auth *Auth) error
}
