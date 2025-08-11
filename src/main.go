package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Login struct {
	ID         uint   `json:"id"`
	UserId     uint   `json:"userId"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	LoginToken string `json:"loginToken"`
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/weblog?"
	charset := "utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn+charset), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Login{})

	AppRun(db)
}
