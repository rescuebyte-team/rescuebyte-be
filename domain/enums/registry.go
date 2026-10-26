package enums

import "reflect"

type EnumDescriptor struct {
	QueryKey    string
	ResponseKey string
	Values      []any
}

var enumRegistry []EnumDescriptor

func RegisterEnum(queryKey, responseKey string, m any) {
	var values []any

	mv := reflect.ValueOf(m)
	if mv.Kind() != reflect.Map {
		return
	}

	for _, k := range mv.MapKeys() {
		v := mv.MapIndex(k).Interface()
		rv := reflect.ValueOf(v)
		keyField := rv.FieldByName("Key")
		if keyField.IsValid() && keyField.String() != "unknown" {
			values = append(values, v)
		}
	}

	enumRegistry = append(enumRegistry, EnumDescriptor{
		QueryKey:    queryKey,
		ResponseKey: responseKey,
		Values:      values,
	})
}

func GetEnumRegistry() []EnumDescriptor {
	return enumRegistry
}
