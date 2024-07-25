package collections

import (
	"fmt"
	"rsdb/src/util"
	attributeDefinitions "rsdb/src/util/attrDefinitions"
	"rsdb/src/util/types"
)

type creatCollectionRequest struct {
	Name   string                 `json:"name" binding:"required"`
	Schema map[string]interface{} `json:"schema"`
}

type createCollectionResponse struct {
	util.Response
	CollectionName string `json:"collection_name"`
}

func getKeys[T any](input map[string]T) []string {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

func checkAttributes(attr map[string]interface{}, schemaKey string) (bool, string) {
	hasTypeKey := false
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
				// TODO: Check Content Attribute
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
	}
	if !hasTypeKey {
		return false, "Key `type` was not found in schema, which is needed to define Attribute types."
	}
	return true, ""
}

func isValidSchema(schema map[string]interface{}) (bool, string) {
	// Loop through keys in schema
	for schemaKey, schemaValue := range schema {
		if schemaKey == "id" {
			return false, "Key `id` was found in document, which is system reserved."
		}
		valid, message := checkAttributes(schemaValue.(map[string]interface{}), schemaKey)
		if !valid {
			return false, message
		}
	}

	return true, ""
}
