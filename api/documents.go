package api

import (
	"encoding/json"
)

type Document struct {
	BaseModel

	ID string `json:"id"`
}

type ListDocumentsResponse struct {
	Count     int               `json:"_count"`
	Documents []json.RawMessage `json:"Documents"`
}
