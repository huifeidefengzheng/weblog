package main

import "gorm.io/gorm"

type User struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
