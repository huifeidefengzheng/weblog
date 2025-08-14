package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AppRun(db *gorm.DB) {
	router := gin.Default()
	pubilcApi := router.Group("/user")
	{
		pubilcApi.GET("/query/all", func(c *gin.Context) {
			QueryAll(db, c)
		})

		pubilcApi.POST("/create", func(c *gin.Context) {
			CreateUser(db, c)
		})
		pubilcApi.POST("/login", func(c *gin.Context) {
			LoginUser(c, db)
		})

	}
	authUserApi := router.Group("/user")
	authUserApi.Use(AuthMiddleware())
	{
		authUserApi.GET("/query/:id", func(c *gin.Context) {
			QueryUserById(db, c)
		})
		authUserApi.POST("/update", func(c *gin.Context) {
			UpdateUser(db, c)
		})
		authUserApi.POST("/delete", func(c *gin.Context) {
			DeleteUser(db, c)
		})

	}

	PostsApi := router.Group("/posts")
	PostsApi.Use(AuthMiddleware())
	{
		PostsApi.POST("/create", func(c *gin.Context) {
			CreatePosts(db, c)
		})
		PostsApi.GET("/query/all", func(c *gin.Context) {
			QueryAllPosts(db, c)
		})
		PostsApi.POST("/update", func(c *gin.Context) {
			UpdatePosts(db, c)
		})

		PostsApi.POST("/delete/:id", func(c *gin.Context) {
			DeletePosts(db, c)
		})
		PostsApi.POST("/:id/comments", func(c *gin.Context) {
			CreateComment(db, c)
		})
	}

	router.Run(":8080")
}
