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

type ItemRepository interface {
	CreateItem(request CreateItemRequest) (ItemRow, error)
	UpdateItem(request UpdateItemRequest) (ItemRow, error)
	GetItem(id uuid.UUID) (ItemRow, error)
	GetItems() ([]ItemRow, error)
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
	err := r.db.Get(&item, `INSERT INTO items (name, description, unit_price, created_by, last_changed_by) VALUES ($1, $2, $3, $4, $4) returning *`, request.Name, request.Description, request.UnitPrice, request.CreatedBy)
	return item, err
}

func (r *ItemRepositoryImpl) UpdateItem(request UpdateItemRequest) (ItemRow, error) {
	var item ItemRow
	err := r.db.Get(&item, `UPDATE items SET name = $1, description = $2, unit_price = $3, last_changed_by = $4 WHERE alt_id = $5 returning *`, request.Name, request.Description, request.UnitPrice, request.LastChangedBy, request.Id)
	return item, err
}
