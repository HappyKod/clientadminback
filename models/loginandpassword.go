package models

type LoginPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}