package clientadminback

type ClientInfo struct {
	TokenExpires string `json:"token_expires"`
	Created      string `json:"created"`
	Id           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	IsActive     bool   `json:"is_active"`
}
