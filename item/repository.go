package item

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"inventory-service-go/commons"
)

type ItemRow struct {
	Id            int64     `db:"id"`
	AltId         uuid.UUID `db:"alt_id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	UnitPrice     float64   `db:"unit_price"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     string    `db:"created_at"`
	LastChangedBy string    `db:"last_changed_by"`
	LastUpdate    string    `db:"last_update"`
}

type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	CreatedBy   string  `json:"created_by"`
}

type UpdateItemRequest struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	UnitPrice     float64   `json:"unit_price"`
	LastChangedBy string    `json:"last_changed_by"`
}

const (
	CREATE_STATEMENT              = "INSERT INTO items (name, description, unit_price, created_by, last_changed_by) VALUES ($1, $2, $3, $4, $4) returning *"
	UPDATE_STATEMENT              = "UPDATE items SET name = $1, description = $2, unit_price = $3, last_changed_by = $4 WHERE alt_id = $5 returning *"
	GET_BY_ID_QUERY               = "SELECT * FROM items WHERE alt_id = $1"
	GET_ALL_QUERY                 = "SELECT * FROM items"
	GET_ALL_QUERY_WITH_PAGINATION = "SELECT * FROM items WHERE id > $1 LIMIT $2"
	DELETE_BY_ID_QUERY            = "DELETE FROM items WHERE alt_id = $1"
)

type ItemRepository interface {
	CreateItem(request CreateItemRequest) (ItemRow, error)
	UpdateItem(request UpdateItemRequest) (ItemRow, error)
	GetItem(id uuid.UUID) (ItemRow, error)
	GetItems(pagination *commons.Pagination) ([]ItemRow, error)
	DeleteItem(id uuid.UUID) (commons.DeleteResult, error)
}

type ItemRepositoryImpl struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepositoryImpl {
	return &ItemRepositoryImpl{
		db: db,
	}
}

func (r *ItemRepositoryImpl) CreateItem(request CreateItemRequest) (ItemRow, error) {
	var item ItemRow
	err := r.db.Get(&item, CREATE_STATEMENT, request.Name, request.Description, request.UnitPrice, request.CreatedBy)
	return item, err
}

func (r *ItemRepositoryImpl) UpdateItem(request UpdateItemRequest) (ItemRow, error) {
	var item ItemRow
	err := r.db.Get(&item, UPDATE_STATEMENT, request.Name, request.Description, request.UnitPrice, request.LastChangedBy, request.Id)
	return item, err
}

func (r *ItemRepositoryImpl) GetItem(id uuid.UUID) (ItemRow, error) {
	var item ItemRow
	err := r.db.Get(&item, GET_BY_ID_QUERY, id)
	return item, err
}

func (r *ItemRepositoryImpl) GetItems(pagination *commons.Pagination) ([]ItemRow, error) {
	var items []ItemRow
	if pagination == nil {
		err := r.db.Select(&items, GET_ALL_QUERY)
		return items, err
	} else {
		err := r.db.Select(&items, GET_ALL_QUERY_WITH_PAGINATION, pagination.LastId, pagination.PageSize)
		return items, err
	}
}

func (r *ItemRepositoryImpl) DeleteItem(id uuid.UUID) (commons.DeleteResult, error) {
	var result commons.DeleteResult
	err := r.db.Get(&result, DELETE_BY_ID_QUERY, id)
	return result, err
}
