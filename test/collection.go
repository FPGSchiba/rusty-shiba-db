package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	router2 "rsdb/src/router"
	"strings"
	"testing"
	"time"
)

func ShouldCreateCollection201(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	requestBody := map[string]interface{}{
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
	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	shouldBody := map[string]interface{}{
		"collection_name": "books",
		"status":          "success",
		"message":         "Successfully created collection: `books`",
	}
	assert.Equal(t, shouldBody, responseBody)
}

func ShouldCreateCollection409(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	requestBody := map[string]interface{}{
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
	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func ShouldFailCreateCollection400(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	requestBody := map[string]interface{}{
		"name":   123,
		"schema": "title:string",
	}

	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
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
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/collections/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "books", responseBody["collection_name"])
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, IsValidUUID(responseBody["collection_id"].(string)), true)
	assert.Equal(t, IsValidDate(responseBody["created_at"].(string)), true)
}

func ShouldReadCollection404(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/collections/movies", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func ShouldUpdateCollection200(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	requestBody := map[string]interface{}{
		"schema": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "str",
			},
			"published_year": map[string]interface{}{
				"type": "nbr",
			},
			"author": map[string]interface{}{
				"type": "str",
			},
		},
	}

	requestJSON, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PATCH", "/api/v1/collections/books", strings.NewReader(string(requestJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "books", responseBody["collection_name"])
}

func ShouldDeleteCollection200(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/api/v1/collections/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "success", responseBody["status"])
}

func ShouldDeleteCollection404(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/api/v1/collections/movies", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "error", responseBody["status"])
}

func createCollection(t *testing.T, collectionName string) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	requestBody := map[string]interface{}{
		"name": collectionName,
	}
	requestJSON, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/collections/", strings.NewReader(string(requestJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusCreated)
}

func ShouldListCollectionPagination200(t *testing.T) {
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	for i := range 10 {
		createCollection(t, fmt.Sprintf("test-%d", i))
	}

	req, _ := http.NewRequest("GET", "/api/v1/collections/?limit=5&offset=0", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, len(responseBody["data"].([]interface{})), 5)
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["limit"], float64(5))
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["offset"], float64(0))
	assert.Equal(t, responseBody["pagination"].(map[string]interface{})["total"], float64(10))
}
