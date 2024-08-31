package commons

import "github.com/google/uuid"

// struct that represents results of a delete operation
type DeleteResult struct {
	Id      uuid.UUID `json:"id"`
	Deleted bool      `json:"deleted"`
}
