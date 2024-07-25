package documents

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"rsdb/src/util"
)

func CreateDocument(c *gin.Context) {
	body := documentCreateRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		}
		return
	}
	for key, value := range body.Data {
		if key == "id" {
			c.JSON(http.StatusBadRequest, util.GetResponseWithMessage("Key `id` was found in document, which is system reserved."))
			return
		}
		fmt.Println(fmt.Sprintf("key: %s;value: %v", key, value))
	}
	c.JSON(http.StatusCreated, documentCreateResponse{
		Response:   util.Response{Status: "success", Message: "Document created successfully."},
		DocumentId: "Testing the waters",
	})
}
