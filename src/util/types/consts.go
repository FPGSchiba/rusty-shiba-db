package types

const (
	STRING = "str"
	NUMBER = "nbr"
	BOOL   = "bol"
	NULL   = "nil"
	ARRAY  = "arr"
	OBJECT = "obj"
	DATE   = "dat"
	UUID   = "uid"
)

func InvalidDataType(typ string) bool {
	result := true
	for _, dataTyp := range []string{STRING, NUMBER, BOOL, NULL, ARRAY, OBJECT, DATE, UUID} {
		if dataTyp == typ {
			result = false
		}
	}
	return result
}
