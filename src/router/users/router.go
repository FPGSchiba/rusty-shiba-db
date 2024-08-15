package users

import (
	"github.com/gin-gonic/gin"
	"rsdb/src/util"
)

func CreateUser(c *gin.Context) {
	c.JSON(200, util.Response{Status: "", Message: ""})
}
