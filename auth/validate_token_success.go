package auth

type ValidateTokenSuccess struct {
	ClientId  string   `json:"client_id"`
	ExpiresIn int64    `json:"expires_in"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserId    string   `json:"user_id"`
}
