package main

import (
	"github.com/gin-gonic/gin"
	"rsdb/src/router"
)

func main() {
	// TODO: Config stuff
	gin.SetMode(gin.DebugMode)
	engine := router.GetRouter()
	engine.Run(":3000")
}
