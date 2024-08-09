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

	collection, message := collections.CreateNewCollection(body.Name, body.Schema)
	if collection == nil { // TODO: More status codes for different cases here
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(message))
		return
	}

	c.JSON(http.StatusCreated, createCollectionResponse{
		Response:       util.Response{Status: "success", Message: message},
		CollectionName: collection.Name,
	})
}

func ReadCollection(c *gin.Context) {
	collectionName := c.Param("collection")
	collection, message := collections.ReadCollection(collectionName)
	if collection == nil {
		c.JSON(http.StatusNotFound, util.GetResponseWithMessage(message))
		return
	}
	c.JSON(http.StatusOK, readCollectionResponse{
		Response:       util.Response{Status: "success", Message: message},
		CollectionName: collectionName,
		CollectionId:   collection.Id,
		Schema:         collection.Schema,
		CreatedAt:      collection.CreatedAt,
		UpdatedAt:      collection.UpdatedAt,
	})
}
