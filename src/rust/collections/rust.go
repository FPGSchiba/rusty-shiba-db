package collections

import (
	"errors"
	"fmt"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"path"
	"runtime"
	"strings"
)

type Collections struct {
	rootPath string
}

type Collection struct {
	Name   string
	Schema map[string]interface{}
	id     string
}

const (
	collectionsFileName = "collections.rsc"

	windowsPath = "\\Data\\RSDB\\"
	linuxPath   = "/data/rsdb/"
	macPath     = "/data/rsdb/"
)

func getRootPath() string {
	var rootPath string
	switch runtime.GOOS {
	case "windows":
		ex, err := os.Executable()
		if err != nil {
			panic(err.Error())
		}

		drive := strings.Split(ex, string(os.PathSeparator))[0]
		rootPath = drive + windowsPath
	case "darwin":
		rootPath = macPath
	case "linux":
		rootPath = linuxPath
	default:
		panic("OS not supported")
	}
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		err := os.MkdirAll(rootPath, os.ModeDir)
		if err != nil {
			panic(err.Error())
		}
	}
	return rootPath
}

func writeCollectionsFile(data []map[string]interface{}) error {
	rootPath := getRootPath()
	collectionsFile := path.Join(rootPath, collectionsFileName)
	value, err := bson.Marshal(bson.M{"collections": data})
	if err != nil {
		return err
	}

	err = os.WriteFile(collectionsFile, value, 0777)
	if err != nil {
		return err
	}
	return nil
}

func readCollectionsFile() ([]map[string]interface{}, error) {
	rootPath := getRootPath()
	collectionsFile := path.Join(rootPath, collectionsFileName)
	file, err := os.ReadFile(collectionsFile)
	if err != nil {
		return nil, err
	}
	println(file)
	var collectionsValue map[string]interface{}
	err = bson.Unmarshal(file, collectionsValue)
	if err != nil {
		return nil, err
	}
	existingCollections := collectionsValue["collections"].([]map[string]interface{})
	return existingCollections, nil
}

func writeSchemaFile(collectionId string, schema map[string]interface{}) error {
	rootPath := getRootPath()
	schemaFile := path.Join(rootPath, fmt.Sprintf("%s.rsc", collectionId))
	marshal, err := bson.Marshal(bson.M{"schema": schema})
	if err != nil {
		return err
	}
	err = os.WriteFile(schemaFile, marshal, 0777)
	if err != nil {
		return err
	}
	return nil
}

func writeIndexFile(collectionId string, index []map[string]interface{}) error {
	rootPath := getRootPath()
	schemaFile := path.Join(rootPath, fmt.Sprintf("%s.rsi", collectionId))
	marshal, err := bson.Marshal(bson.M{"index": index})
	if err != nil {
		return err
	}
	err = os.WriteFile(schemaFile, marshal, 0777)
	if err != nil {
		return err
	}
	return nil
}

func InitRustyStorage() *Collections {
	// Root directory for DB exists
	rootPath := getRootPath()
	_, err := os.ReadDir(rootPath)
	if err != nil {
		return nil
	}

	collectionsFile := path.Join(rootPath, collectionsFileName)
	// Check collections file
	if _, err := os.Stat(collectionsFile); errors.Is(err, os.ErrNotExist) {
		// Init collections file
		err := writeCollectionsFile([]map[string]interface{}{})
		if err != nil {
			return nil
		}
	}

	return &Collections{rootPath: rootPath}
}

func CreateNewCollection(name string, schema map[string]interface{}) *Collection {
	collId := uuid.NewUUID().String()

	// Add to collections
	existingCollections, err := readCollectionsFile()
	if err != nil {
		log.WithFields(log.Fields{
			"component":  "RustyStorage",
			"function":   "CreateNewCollection",
			"collection": name,
		}).Error(fmt.Sprintf("Failed to read existing collections: %s", err.Error()))
		return nil // TODO: Correct error handling
	}
	for _, collection := range existingCollections {
		if collection["name"] == name {
			log.WithFields(log.Fields{
				"component":  "RustyStorage",
				"function":   "CreateNewCollection",
				"collection": name,
			}).Error("Collection already exists")
			return nil // TODO: Correct error handling
		}
	}
	newCollectionValue := map[string]interface{}{
		"name": name,
		"id":   collId,
	}
	newCollections := append(existingCollections, newCollectionValue)
	err = writeCollectionsFile(newCollections)
	if err != nil {
		log.WithFields(log.Fields{
			"component":  "RustyStorage",
			"function":   "CreateNewCollection",
			"collection": name,
		}).Error(fmt.Sprintf("Failed to update existing Collections: %s", err.Error()))
		return nil // TODO: Correct error handling
	}

	// Create Schema file
	if schema != nil {
		err := writeSchemaFile(collId, schema)
		if err != nil {
			log.WithFields(log.Fields{
				"component":  "RustyStorage",
				"function":   "CreateNewCollection",
				"collection": name,
			}).Error(fmt.Sprintf("Failed to create Schema file: %s", err.Error()))
			return nil // TODO: Correct error handling
		}
	}

	// Create Index file
	err = writeIndexFile(collId, newCollections)
	if err != nil {
		log.WithFields(log.Fields{
			"component":  "RustyStorage",
			"function":   "CreateNewCollection",
			"collection": name,
		}).Error(fmt.Sprintf("Failed to create Index file: %s", err.Error()))
		return nil // TODO: Correct error handling
	}

	// Create Data folder
	err = os.MkdirAll(path.Join(getRootPath(), collId), os.ModeDir)
	if err != nil {
		log.WithFields(log.Fields{
			"component":  "RustyStorage",
			"function":   "CreateNewCollection",
			"collection": name,
		}).Error(fmt.Sprintf("Failed to create Data folder: %s", err.Error()))
		return nil
	}

	return &Collection{Name: name, Schema: schema, id: collId}
}
