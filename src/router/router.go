package router

import (
	"github.com/gin-gonic/gin"
	"rsdb/src/router/collections"
	"rsdb/src/router/documents"
	"rsdb/src/util"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(util.JSONLogMiddleware())
	router.Use(util.RequestID(util.RequestIDOptions{AllowSetting: false}))
	router.Use(util.CORS(util.CORSOptions{}))

	superGroup := router.Group("/api/v1")
	{
		superGroup.GET("/", getVersion)
		documentGroup := superGroup.Group("/:collection/documents")
		{
			documentGroup.POST("/", documents.CreateDocument)
		}
		collectionGroup := superGroup.Group("/collections")
		{
			collectionGroup.POST("/", collections.CreateCollection)
			collectionGroup.GET("/:collection", collections.ReadCollection)
		}
	}
	return router
}
