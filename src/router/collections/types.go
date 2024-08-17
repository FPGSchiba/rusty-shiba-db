package collections

import (
	"fmt"
	"regexp"
	"rsdb/src/rust/collections"
	"rsdb/src/util"
	attributeDefinitions "rsdb/src/util/attrDefinitions"
	"rsdb/src/util/types"
)

type creatCollectionRequest struct {
	Name   string                 `json:"name" binding:"required" validations:"type=string"`
	Schema map[string]interface{} `json:"schema"`
}

type createCollectionResponse struct {
	util.Response
	CollectionName string `json:"collection_name"`
}

type readCollectionResponse struct {
	util.Response
	CollectionName string                 `json:"collection_name"`
	CollectionId   string                 `json:"collection_id"`
	Schema         map[string]interface{} `json:"schema"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
}

type updateCollectionRequest struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

type updateCollectionResponse struct {
	util.Response
	CollectionName string          `json:"collection_name"`
	Pagination     util.Pagination `json:"pagination"`
}

type listCollectionsResponse struct {
	util.Response
	Data       []collections.CollectionInfo `json:"data"`
	Pagination util.Pagination              `json:"pagination"`
}

func getKeys[T any](input map[string]T) []string {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

func getTypeFromAttr(attr map[string]interface{}) (bool, string) {
	hasType := false
	var dataType string
	for attributeKey, attributeValue := range attr {
		if attributeKey == "type" {
			hasType = true
			dataType = attributeValue.(string)
		}
	}
	if !hasType {
		return false, ""
	}
	return true, dataType
}

func checkAttributes(attr map[string]interface{}, schemaKey string) (bool, string) {
	hasTypeKey, dataType := getTypeFromAttr(attr)
	if !hasTypeKey {
		return false, fmt.Sprintf("Key `type` was not found in schema for Key: `%s`, which is needed to define Attribute types.", schemaKey)
	}
	for attributeKey, attributeValue := range attr {
		if attributeDefinitions.InvalidAttributeDefinition(attributeKey) {
			return false, fmt.Sprintf("Attribute key: `%s` does not exist, schema error.", attributeKey)
		}

		if attributeKey == "type" {
			hasTypeKey = true
			if types.InvalidDataType(attributeValue.(string)) {
				return false, fmt.Sprintf("Datatype: `%s` does not exist, schema error.", attributeValue.(string))
			}

			if attributeValue.(string) == types.ARRAY {
				keys := getKeys(attr)
				itemsExist := false
				for _, key := range keys {
					if key == "items" {
						itemsExist = true
					}
				}

				if !itemsExist {
					return false, fmt.Sprintf("Key `items` does not exist for array: `%s`, schema error.", schemaKey)
				}

				valid, message := checkAttributes(attr["items"].(map[string]interface{}), schemaKey)
				if !valid {
					return false, message
				}
			}

			if attributeValue.(string) == types.OBJECT {
				keys := getKeys(attr)
				contentExists := false
				for _, key := range keys {
					if key == "content" {
						contentExists = true
					}
				}

				if !contentExists {
					return false, fmt.Sprintf("Key `content` does not exist for array: `%s`, schema error.", schemaKey)
				}

				valid, message := isValidSchema(attr["content"].(map[string]interface{}))
				if !valid {
					return false, message
				}
			}
		}

		if attributeDefinitions.InvalidAttributeDefinitionForType(dataType, attributeKey) {
			return false, fmt.Sprintf("Attribute key: `%s` does not exist for dataType: `%s` in schemaKey: `%s`, schema error.", attributeKey, dataType, schemaKey)
		}
	}

	return true, ""
}

func isValidSchema(schema map[string]interface{}) (bool, string) {
	// Loop through keys in schema
	for schemaKey, schemaValue := range schema {
		if schemaKey == "id" {
			return false, "Key `id` was found in document, which is system reserved."
		}
		switch schemaValue.(type) {
		case map[string]interface{}:
			valid, message := checkAttributes(schemaValue.(map[string]interface{}), schemaKey)
			if !valid {
				return false, message
			}
			break
		default:
			return false, fmt.Sprintf("Schema key `%s` is not a valid schema attribute.", schemaKey)
		}
	}

	return true, ""
}

func isValidName(name string) bool {
	compile, err := regexp.Compile("^[a-z0-9-]*$")
	if err != nil {
		return false
	}
	return compile.MatchString(name)
}
