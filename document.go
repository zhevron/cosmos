package cosmos

import (
	"reflect"
)

func DocumentID(document interface{}) (string, error) {
	rv := reflect.ValueOf(document)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", &CosmosError{Code: ErrNoDocumentID, Message: "document is not a struct"}
	}

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

	return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not find id field in struct"}
}
