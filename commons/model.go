package commons

import "github.com/google/uuid"

// struct that represents results of a delete operation
type DeleteResult struct {
	Id      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}

type Pagination struct {
	LastId   int `json:"last_id"`
	PageSize int `json:"page_size"`
}

type AuditInfo struct {
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	LastUpdate    string `json:"last_update"`
	LastChangedBy string `json:"last_change_by"`
}
