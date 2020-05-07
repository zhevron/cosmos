package cosmos

import (
	"reflect"
	"strings"
)

func DocumentID(document interface{}) (string, error) {
	if doc, ok := document.(Document); ok {
		return doc.ID, nil
	}

	rv := reflect.ValueOf(document)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		return extractDocumentIDFromMap(rv)

	case reflect.Struct:
		return extractDocumentIDFromStruct(rv)
	}

	return "", &CosmosError{Code: ErrNoDocumentID, Message: "unsupported document type"}
}

func extractDocumentIDFromMap(document reflect.Value) (string, error) {
	keys := document.MapKeys()
	if len(keys) == 0 {
		return "", &CosmosError{Code: ErrNoDocumentID, Message: "map has no keys"}
	}

	if keys[0].Kind() != reflect.String {
		return "", &CosmosError{Code: ErrNoDocumentID, Message: "map keys are not strings"}
	}

	for _, k := range keys {
		if k.String() == "id" {
			if id, ok := document.MapIndex(k).Interface().(string); ok && id != "" {
				return id, nil
			}

			return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not convert id to string"}
		}
	}

	return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not find id key in map"}
}

func extractDocumentIDFromStruct(document reflect.Value) (string, error) {
	rt := document.Type()
	numField := rt.NumField()

	for i := 0; i < numField; i++ {
		if rt.Field(i).Tag.Get("json") == "id" {
			if id, ok := document.Field(i).Interface().(string); ok && id != "" {
				return id, nil
			}

			return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not convert id to string"}
		}
	}

	for i := 0; i < numField; i++ {
		if strings.EqualFold(rt.Field(i).Name, "id") {
			if id, ok := document.Field(i).Interface().(string); ok && id != "" {
				return id, nil
			}

			return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not convert id to string"}
		}
	}

	return "", &CosmosError{Code: ErrNoDocumentID, Message: "could not find id field in struct"}
}
