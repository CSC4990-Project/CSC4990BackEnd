package models

type User struct {
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

type Usertype struct {
	UserType string `json:"usertype"`
	ID       int    `json:"id" gorm:"primary_key"`
}
type EmailType struct {
	Email    string `json:"email"`
	UserType string `json:"userType"`
}
