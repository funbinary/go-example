package main

import (
	"github.com/bin-work/go-example/pkg/bfile"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	e := gin.Default()

	e.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"key": "value"})
	})

	e.RunTLS(":8443",
		bfile.Join(bfile.SelfDir(), "cert/server.pem"),
		bfile.Join(bfile.SelfDir(), "cert/server.key"))
}
