package cosmos

import (
	"reflect"
)

func DocumentID(document interface{}) (string, error) {
	rv := reflect.ValueOf(document)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		keys := rv.MapKeys()
		if len(keys) == 0 {
			return "", &CosmosError{Code: ErrNoDocumentID, Message: "map has no keys"}
		}

		if keys[0].Kind() != reflect.String {
			return "", &CosmosError{Code: ErrNoDocumentID, Message: "map keys are not strings"}
		}

		for _, k := range keys {
			if k.String() == "id" {
				if id, ok := rv.MapIndex(k).Interface().(string); ok && id != "" {
					return id, nil
				}

				return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not convert id to string"}
			}
		}

	case reflect.Struct:
		rt := rv.Type()
		numField := rt.NumField()
		for i := 0; i < numField; i++ {
			if rt.Field(i).Tag.Get("json") == "id" {
				if id, ok := rv.Field(i).Interface().(string); ok && id != "" {
					return id, nil
				}

				return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not convert id to string"}
			}
		}
	}

	return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not find id field in struct"}
}
