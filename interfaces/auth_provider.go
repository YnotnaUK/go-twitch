package interfaces

type AuthProvider interface {
	GetAccessToken() (string, error)
	GetLoginAndAccessToken() (string, string, error)
}
