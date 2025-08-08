package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}

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

	router := gin.Default()
	router.GET("/user/query/all", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(200, gin.H{
			"data":    users,
			"status":  "success",
			"message": "success",
		})
	})

	router.POST("/user/create", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Table("users").Where("name = ?", json.Name).First(&json)
		if json.ID != 0 {
			c.JSON(400, gin.H{
				"status":  "用户名称已存在",
				"message": json.Name,
			})
			return
		}
		var email = json.Email
		var user = User{
			Name:     json.Name,
			Password: json.Password,
			Email:    email,
		}
		db.Create(&user)
		c.JSON(200, gin.H{
			"data":    user.ID,
			"status":  "success",
			"message": "success",
		})
	})

	// 更新用户信息
	router.POST("/user/update", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return

		}
		if json.ID == 0 {
			c.JSON(400, gin.H{
				"status":  "用户ID不能为空",
				"message": "用户ID不能为空",
			})
			return
		}
		var user User
		db.Table("users").Where("id = ?", json.ID).First(&user)
		if user.ID == 0 {
			c.JSON(400, gin.H{
				"status":  "用户不存在",
				"message": "用户不存在",
			})
			return
		}
		user.Name = json.Name
		user.Password = json.Password
		user.Email = json.Email
		db.Table("users").Where("id = ?", json.ID).Updates(&user)
		c.JSON(200, gin.H{
			"data":    user.ID,
			"status":  "success",
			"message": "success",
		})
	})
	router.Run(":8080")
}
