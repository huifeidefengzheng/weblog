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
		pubilcApi.POST("/update", func(c *gin.Context) {
			UpdateUser(db, c)
		})

		pubilcApi.POST("/delete", func(c *gin.Context) {
			DeleteUser(db, c)
		})
	}

	router.Run(":8080")
}
