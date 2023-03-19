package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("admin:csc4990db@tcp(csc4990db.c3exsdfmiwh2.us-east-2.rds.amazonaws.com:3306)/csc4990?charset=utf8mb4"), &gorm.Config{})
	if err != nil {
		panic("could not connect to the database")
	}
	DB = connection

}
