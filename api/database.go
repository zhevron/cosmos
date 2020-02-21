package api

type Database struct {
	BaseModel

	ID string `json:"id"`
}

type ListDatabasesResponse struct {
	Databases []Database `json:"databases"`
}

// TODO: CreateDatabaseRequest (https://docs.microsoft.com/en-us/rest/api/cosmos-db/create-a-database)
