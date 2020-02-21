package api

type Database struct {
	BaseModel

	ID string `json:"id"`
}

type ListDatabasesResponse struct {
	Databases []Database `json:"databases"`
}

type CreateDatabaseRequest struct {
	ID string `json:"id"`
}
