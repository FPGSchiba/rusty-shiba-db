package attributeDefinitions

import "rsdb/src/util/types"

const (
	TYPE    = "type"
	UNIQUE  = "unique"
	ITEMS   = "items"
	CONTENT = "content"
)

func InvalidAttributeDefinition(key string) bool {
	result := true
	for _, attrDef := range []string{TYPE, UNIQUE, ITEMS, CONTENT} {
		if attrDef == key {
			result = false
		}
	}
	return result
}

func InvalidAttributeDefinitionForType(dataType string, key string) bool {
	result := true
	if InvalidAttributeDefinition(key) {
		return true
	}
	switch dataType {
	case types.STRING:
		for _, attrDef := range []string{TYPE, UNIQUE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.NUMBER:
		for _, attrDef := range []string{TYPE, UNIQUE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.BOOL:
		for _, attrDef := range []string{TYPE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.NULL:
		for _, attrDef := range []string{TYPE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.ARRAY:
		for _, attrDef := range []string{TYPE, UNIQUE, ITEMS} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.OBJECT:
		for _, attrDef := range []string{TYPE, UNIQUE, CONTENT} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.DATE:
		for _, attrDef := range []string{TYPE, UNIQUE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	case types.UUID:
		for _, attrDef := range []string{TYPE, UNIQUE} {
			if attrDef == key {
				result = false
			}
		}
		return result
	}
	return result
}
