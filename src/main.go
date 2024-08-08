package main

import (
	"github.com/gin-gonic/gin"
	"rsdb/src/router"
	"rsdb/src/rust/collections"
)

func main() {
	// TODO: Config stuff

	collections.InitRustyStorage()
	gin.SetMode(gin.DebugMode)
	engine := router.GetRouter()
	err := engine.Run(":3000")
	if err != nil {
		return
	}
}
