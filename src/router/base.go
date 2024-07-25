package router

import (
	"github.com/gin-gonic/gin"
	"rsdb/src/util"
)

func getVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": util.ApiVersion,
	})
}
