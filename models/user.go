package models

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` //never expose the hash in JSON
}
