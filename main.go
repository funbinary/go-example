package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.StaticFS("/", http.Dir("./dist"))
	r.Run(":8081")
}
