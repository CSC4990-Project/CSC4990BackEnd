package models

type User struct {
	UserType uint   `json:"-"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}
