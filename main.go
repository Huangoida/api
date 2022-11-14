package main

import (
	"api/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.NoRoute(handler.DoProxy)
	r.Run(":8082")
}
