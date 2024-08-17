package test

import (
	"github.com/go-playground/assert/v2"
	"github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	router2 "rsdb/src/router"
	"rsdb/src/rust/collections"
	"rsdb/src/util"
	"testing"
)

func TestAPI(t *testing.T) {
	InitDatabase()

	t.Run("FullAPI", func(t *testing.T) {
		t.Run("Collections", func(t *testing.T) {
			t.Parallel()
			t.Run("Create a New Collection Successfully", ShouldCreateCollection201)
			t.Run("Create Collection with Existing Name", ShouldCreateCollection409)
			t.Run("Retrieve a Collection by Name", ShouldReadCollection200)
			t.Run("Retrieve a Non-Existent Collection", ShouldReadCollection404)
			t.Run("Update an Existing Collection", ShouldUpdateCollection200)
			t.Run("Delete a Collection Successfully", ShouldDeleteCollection200)
			t.Run("Delete a Non-Existent Collection", ShouldDeleteCollection404)
			t.Run("List Collections with Pagination", ShouldListCollectionPagination200)
			t.Run("Invalid Request Format During Collection Creation", ShouldFailCreateCollection400)
		})
		t.Run("Base", func(t *testing.T) {
			t.Parallel()
			t.Run("Get Version Successfully", ShouldGetVersion)
		})
	})

	CleanupDatabase()
}

func ShouldGetVersion(t *testing.T) {
	t.Parallel()
	router := router2.GetRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	var body map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &body)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, util.ApiVersion, body["version"].(string))
}

func InitDatabase() {
	collections.InitRustyStorage()
}

func CleanupDatabase() {
	err := collections.DestroyRustyStorage()
	if err != nil {
		panic(err)
	}
}
