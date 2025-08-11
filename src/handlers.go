package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func QueryAll(db *gorm.DB, c *gin.Context) {
	// 创建页面 start limit map指针
	var page = struct {
		Start int `form:"start"`
		Limit int `form:"limit"`
	}{}
	if err := c.ShouldBind(&page); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Println("start:", page.Start, "limit:", page.Limit)
	var users []User
	db.Debug().Table("users").Limit(page.Limit).Offset(page.Start).Find(&users)
	c.JSON(200, gin.H{
		"data":    users,
		"status":  "success",
		"message": "success",
	})
}

func CreateUser(db *gorm.DB, c *gin.Context) {
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
}

// 更新用户信息
func UpdateUser(db *gorm.DB, c *gin.Context) {
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
}

func DeleteUser(db *gorm.DB, c *gin.Context) {
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
	db.Table("users").Where("id = ?", json.ID).Delete(&user)
	c.JSON(200, gin.H{
		"data":    user.ID,
		"status":  "success",
		"message": "delete user success",
	})
}
