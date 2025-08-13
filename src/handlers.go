package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	// 将user的password改成空
	for i := range users {
		users[i].Password = ""
	}
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
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"error": "密码加密失败"})
		return
	}
	var existUser User
	db.Table("users").Where("username = ?", json.Username).First(&existUser)
	if existUser.ID != 0 {
		c.JSON(400, gin.H{
			"status":  "用户名已存在",
			"message": json.Username,
		})
		return
	}
	var email = json.Email
	var user = User{
		Username: json.Username,
		Password: string(hashedPassword),
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
			"status":  "fail",
			"message": "用户不存在",
		})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	hashedPasswordstr := string(hashedPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": "密码加密失败"})
		return
	}

	user.Username = json.Username
	user.Password = hashedPasswordstr
	user.Email = json.Email
	db.Table("users").Where("id = ?", json.ID).Updates(&user)
	c.JSON(200, gin.H{
		"data":    user.ID,
		"status":  "success",
		"message": "update user success",
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

/*func CreatePosts(db *gorm.DB, c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	json.UserID = userID
	var user User

}*/
