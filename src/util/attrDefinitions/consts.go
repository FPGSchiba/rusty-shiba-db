package attributeDefinitions

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
