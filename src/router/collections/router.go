package collections

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"rsdb/src/rust/collections"
	"rsdb/src/util"
	"strconv"
	"strings"
)

func CreateCollection(c *gin.Context) {
	body := creatCollectionRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}

	if !isValidName(body.Name) {
		c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(fmt.Sprintf("Invalid collection name: `%s`, must match regex Pattern: `^[a-z0-9-]*$`", body.Name)))
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
		if strings.Contains(message, "already exists") {
			c.JSON(http.StatusConflict, util.GetResponseWithMessage(message))
			return
		}
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
	if !isValidName(collectionName) {
		c.JSON(http.StatusNotFound, util.GetResponseWithMessage(fmt.Sprintf("Invalid collection name: `%s`, must match regex Pattern: `^[a-z0-9-]*$`", collectionName)))
		return
	}
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

func UpdateCollection(c *gin.Context) {
	collectionName := c.Param("collection")
	body := updateCollectionRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		}
		return
	}
	if !isValidName(body.Name) {
		c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(fmt.Sprintf("Invalid collection name: `%s`, must match regex Pattern: `^[a-z0-9-]*$`", collectionName)))
		return
	}

	if body.Name != "" {
		if !isValidName(body.Name) {
			c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(fmt.Sprintf("Invalid new collection name: `%s`, must match regex Pattern: `^[a-z0-9-]*$`", collectionName)))
			return
		}
		if collectionName == body.Name {
			c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(fmt.Sprintf("Cannot rename collection to same name.")))
			return
		}
		success, message := collections.UpdateCollectionName(collectionName, body.Name)
		if !success {
			c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(message))
			return
		}
		collectionName = body.Name
	}
	if body.Schema != nil {
		isValid, message := isValidSchema(body.Schema)
		if !isValid {
			c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(message))
			return
		}
		success, message := collections.UpdateCollectionSchema(collectionName, body.Schema)
		if !success { // TODO: More status codes
			c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(message))
			return
		}
	}
	c.JSON(http.StatusOK, updateCollectionResponse{
		Response:       util.Response{Status: "success", Message: "Collection updated successfully."},
		CollectionName: collectionName,
	})
}

func DeleteCollection(c *gin.Context) {
	collectionName := c.Param("collection")
	if !isValidName(collectionName) {
		c.JSON(http.StatusNotFound, util.GetResponseWithMessage(fmt.Sprintf("Invalid collection name: `%s`, must match regex Pattern: `^[a-z0-9-]*$`", collectionName)))
		return
	}
	success, message := collections.DeleteCollectionByName(collectionName)
	if !success {
		if strings.Contains(message, "not found") {
			c.JSON(http.StatusNotFound, util.GetResponseWithMessage(message))
			return
		}
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(message))
		return
	}
	c.JSON(http.StatusOK, util.Response{Status: "success", Message: "Collection deleted successfully."})
}

func getLimitAndOffset(c *gin.Context) (int, int, error, string) {
	parLimit, ok := c.GetQuery("limit")
	if !ok {
		parLimit = "10"
	}
	parOffset, ok := c.GetQuery("parOffset")
	if !ok {
		parOffset = "0"
	}
	limit, err := strconv.Atoi(parLimit)
	if err != nil {
		return 0, 0, err, fmt.Sprintf("Failed to parse limit: %s. Needs to be an Integer.", parLimit)
	}
	offset, err := strconv.Atoi(parOffset)
	if err != nil {
		return 0, 0, err, fmt.Sprintf("Failed to parse offset: %s. Needs to be an Integer.", parOffset)
	}
	return limit, offset, nil, ""
}

func paginate(colls []collections.CollectionInfo, offset int, limit int) []collections.CollectionInfo {
	total := len(colls)
	// Validate the offset and limit
	if offset < 0 || offset >= total {
		return []collections.CollectionInfo{} // Return an empty slice if the offset is out of bounds
	}

	if limit <= 0 {
		return []collections.CollectionInfo{} // Return an empty slice if the limit is not positive
	}

	// Calculate the end index for the slice
	end := offset + limit
	if end > total {
		end = total // Ensure the end index does not exceed the total length
	}

	// Return the paginated slice
	return colls[offset:end]
}

func ListCollections(c *gin.Context) {
	limit, offset, err, message := getLimitAndOffset(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetResponseWithMessage(message))
		return
	}
	colls, message := collections.ListAllCollections() // Inefficient, please fix
	if colls == nil {
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(message))
		return
	}
	total := len(colls)
	collectionList := paginate(colls, offset, limit)
	c.JSON(http.StatusOK, listCollectionsResponse{
		Response: util.Response{Status: "success", Message: "Successfully listed all collections."},
		Data:     collectionList,
		Pagination: util.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}
