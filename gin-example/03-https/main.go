package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func main() {
	router := gin.Default()
	//router.Use(TlsHandler())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.RunTLS(":80", "ssl.pem", "ssl.key")
	//router.Run(":8080")
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":80",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
