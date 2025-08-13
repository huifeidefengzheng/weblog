package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥
var jwtKey = []byte("secret_key")
var defalutJWTKey = []byte("default_key")

func init() {
	if len(jwtKey) == 0 {
		jwtKey = defalutJWTKey
	}
}

// 登录凭证
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// JWt 认证
type JwtToken struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LoginUser(c *gin.Context, db *gorm.DB) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid credentials"})
		return
	}
	var user User
	db.Table("users").Where("username = ?", credentials.Username).First(&user)
	if user.ID == 0 {
		c.JSON(400, gin.H{"error": "用户名或者密码错误"})
		return
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(400, gin.H{"error": "用户名或者密码错误"})
		return
	}
	// 生成JWT
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &JwtToken{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(200, gin.H{
		"data":    map[string]string{"token": tokenString, "user_id": strconv.Itoa(int(user.ID)), "username": user.Username},
		"status":  "success",
		"message": "登录成功",
	})

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokening := c.GetHeader("Authorization")
		if tokening == "" {
			c.JSON(401, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}
		var jwtToken = &JwtToken{}
		token, err := jwt.ParseWithClaims(tokening, jwtToken, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 验证成功，将用户ID保存到上下文中
		c.Set("user_id", jwtToken.UserID)
		c.Set("username", jwtToken.Username)
		c.Next()
	}
}
