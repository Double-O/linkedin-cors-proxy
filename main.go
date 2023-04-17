package main

import (
	"github.com/Double-O/linkedin-cors-proxy/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// enabling cors
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/linkedin/v2/me", handlers.HandleLinkedInMe())
	r.GET("/linkedin/v2/accessToken", handlers.HandleLinkedinAccessToken())
	r.POST("/linkedin/v2/ugcPost", handlers.HandleLinkedinPostShare())
	r.POST("/linkedin/v2/assets/registerUpload", handlers.HandleLinkedinRegisterUpload())
	r.POST("/linkedin/v2/assets/upload-file", handlers.HandleLinkedinMediaUpload())

	r.Run() // listen and serve on 0.0.0.0:8080

}
