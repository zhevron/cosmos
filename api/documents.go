package api

import (
	"encoding/json"
)

type Document struct {
	BaseModel

	ID string `json:"id"`
}

type QueryDocumentsResponse struct {
	Count     int64             `json:"_count"`
	Documents []json.RawMessage `json:"Documents"`
}
