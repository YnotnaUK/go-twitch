package entities

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
