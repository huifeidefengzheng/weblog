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

// 根据ID查询用户信息
func QueryUserById(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")

	var user User
	db.Table("users").Where("id = ?", id).First(&user)
	if user.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "用户不存在",
			"message": "用户不存在",
		})
		return
	}
	user.Password = ""
	c.JSON(200, gin.H{
		"data":    user,
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

func CreatePosts(db *gorm.DB, c *gin.Context) {
	log.Println(c)
	userID := c.MustGet("user_id").(uint)
	log.Println("userID:", userID)
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	json.UserID = userID
	if err := db.Create(&json).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "创建文章失败",
		})
	}
	c.JSON(200, gin.H{
		"data":    json,
		"status":  "success",
		"message": "创建文章成功",
	})
}

func QueryAllPosts(db *gorm.DB, c *gin.Context) {
	log.Println(c)
	userID := c.MustGet("user_id").(uint)
	log.Println("userID:", userID)
	var posts []Post
	res := db.Preload("User").Preload("Comments").First(&posts)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "查询文章失败",
		})
	}
	// 将所有用户的密码改成空
	for i := range posts {
		posts[i].User.Password = ""
	}
	c.JSON(200, gin.H{
		"data":    posts,
		"status":  "success",
		"message": "查询文章成功",
	})
}

// UpdatePosts 更新文章
func UpdatePosts(db *gorm.DB, c *gin.Context) {
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if json.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "文章ID不能为空",
			"message": "文章ID不能为空",
		})
		return
	}
	var post Post
	db.Table("posts").Where("id = ?", json.ID).First(&post)
	if post.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "文章不存在",
			"message": "文章不存在",
		})
	}
	post.Title = json.Title
	post.Content = json.Content
	db.Table("posts").Where("id = ?", json.ID).Updates(&post)
	c.JSON(200, gin.H{
		"data":    post,
		"status":  "success",
		"message": "更新文章成功",
	})
}

// DeletePosts 删除文章
func DeletePosts(db *gorm.DB, c *gin.Context) {
	id := c.Param("id")
	log.Println("id:", id)
	var post Post
	db.Table("posts").Where("id = ?", id).First(&post)
	if post.ID == 0 {
		c.JSON(400, gin.H{
			"status":  "文章不存在",
			"message": "文章不存在",
			"data":    post.ID,
		})
	}
	db.Table("posts").Where("id = ?", id).Delete(&post)
	c.JSON(200, gin.H{
		"data":    post.ID,
		"status":  "success",
		"message": "删除文章成功",
	})
}

// CreateComment 创建评论
func CreateComment(db *gorm.DB, c *gin.Context) {
	post_id := c.Param("id")
	log.Println(post_id)
	userID := c.MustGet("user_id").(uint)
	log.Println("userID:", userID)
	var json Comment
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var common Comment
	common.Content = json.Content
	common.PostID = json.PostID
	common.UserID = userID
	if err := db.Create(&common).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "创建评论失败",
		})
	}
	c.JSON(200, gin.H{
		"data":    common,
		"status":  "success",
		"message": "创建评论成功",
	})

}
