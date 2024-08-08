package collections

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"rsdb/src/rust/collections"
	"rsdb/src/util"
)

func CreateCollection(c *gin.Context) {
	body := creatCollectionRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		}
		return
	}
	if body.Schema != nil {
		valid, message := isValidSchema(body.Schema)
		if !valid {
			c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(message))
			return
		}
	}

	collection := collections.CreateNewCollection(body.Name, body.Schema)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage("collection creation failed"))
		return
	}

	c.JSON(http.StatusCreated, createCollectionResponse{
		Response:       util.Response{Status: "success", Message: "Collection created successfully."},
		CollectionName: collection.Name,
	})
}
