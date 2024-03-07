package auth

type ValidateTokenFailed struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
