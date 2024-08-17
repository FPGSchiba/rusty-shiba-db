package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"testing"
	"time"
)

func getBooksSchema() map[string]interface{} {
	return map[string]interface{}{
		"name": "books",
		"schema": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "str",
			},
			"published_year": map[string]interface{}{
				"type": "nbr",
			},
		},
	}
}

func ShouldCreateCollection201(t *testing.T) {
	requestBody := getBooksSchema()
	requestJSON, _ := json.Marshal(requestBody)
	resp, err := request("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}

	shouldBody := map[string]interface{}{
		"collection_name": "books",
		"status":          "success",
		"message":         "Successfully created collection: `books`",
	}
	assert.Equal(t, shouldBody, responseBody)
}

func ShouldCreateCollection409(t *testing.T) {
	requestBody := getBooksSchema()
	requestJSON, _ := json.Marshal(requestBody)
	resp, err := request("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func ShouldFailCreateCollection400(t *testing.T) {
	requestBody := map[string]interface{}{
		"name":   123,
		"schema": "title:string",
	}
	requestJSON, _ := json.Marshal(requestBody)
	resp, err := request("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidDate(date string) bool {
	_, err := time.Parse(time.RFC3339, date)
	return err == nil
}

func ShouldReadCollection200(t *testing.T) {
	resp, err := request("GET", "/api/v1/collections/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "books", responseBody["collection_name"])
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, IsValidUUID(responseBody["collection_id"].(string)), true)
	assert.Equal(t, IsValidDate(responseBody["created_at"].(string)), true)
}

func ShouldReadCollection404(t *testing.T) {
	resp, err := request("GET", "/api/v1/collections/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func ShouldUpdateCollection200(t *testing.T) {
	requestBody := getBooksSchema()
	requestBody["author"] = map[string]interface{}{
		"type": "str",
	}
	requestJSON, _ := json.Marshal(requestBody)
	resp, err := request("PATCH", "/api/v1/collections/books", strings.NewReader(string(requestJSON)))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "books", responseBody["collection_name"])
}

func ShouldDeleteCollection200(t *testing.T) {
	resp, err := request("DELETE", "/api/v1/collections/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	assert.Equal(t, "success", responseBody["status"])
}

func ShouldDeleteCollection404(t *testing.T) {
	resp, err := request("DELETE", "/api/v1/collections/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func createCollection(t *testing.T, collectionName string) {
	requestBody := map[string]interface{}{
		"name": collectionName,
	}
	requestJSON, _ := json.Marshal(requestBody)
	resp, err := request("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
}

func ShouldListCollectionPagination200(t *testing.T) {
	for i := range 10 {
		createCollection(t, fmt.Sprintf("test-%d", i))
	}

	resp, err := request("GET", "/api/v1/collections/?limit=5&offset=5", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	responseBody, err := getResponseBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, len(responseBody["data"].([]interface{})), 5)
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["limit"], float64(5))
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["offset"], float64(0))
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["total"], float64(10))
}
