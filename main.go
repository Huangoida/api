package main

import (
	"api/Init"
	"api/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	Init.InitConfig()
	r := gin.Default()

	r.NoRoute(handler.DoProxy)
	r.Run(":8082")
}
