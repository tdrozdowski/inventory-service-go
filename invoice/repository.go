package invoice

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"inventory-service-go/commons"
	"time"
)

type InvoiceRow struct {
	Id            int64     `db:"id"`
	AltId         uuid.UUID `db:"alt_id"`
	UserId        uuid.UUID `db:"user_id"`
	Total         float64   `db:"total"`
	Paid          bool      `db:"paid"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	LastChangedBy string    `db:"last_changed_by"`
	LastUpdate    time.Time `db:"last_update"`
}

type InvoiceItemRow struct {
	Id                int64           `db:"id"`
	AltId             uuid.UUID       `db:"alt_id"`
	UserId            uuid.UUID       `db:"user_id"`
	Total             float64         `db:"total"`
	Paid              bool            `db:"paid"`
	CreatedBy         string          `db:"created_by"`
	CreatedAt         time.Time       `db:"created_at"`
	LastChangedBy     string          `db:"last_changed_by"`
	LastUpdate        time.Time       `db:"last_update"`
	ItemSeqId         sql.NullInt64   `db:"item_seq"`
	ItemAltId         uuid.UUID       `db:"item_alt_id"`
	ItemName          sql.NullString  `db:"item_name"`
	ItemDescription   sql.NullString  `db:"item_description"`
	ItemUnitPrice     sql.NullFloat64 `db:"item_unit_price"`
	ItemCreatedBy     sql.NullString  `db:"item_created_by"`
	ItemCreatedAt     sql.NullTime    `db:"item_created_at"`
	ItemLastChangedBy sql.NullString  `db:"item_last_changed_by"`
	ItemLastUpdate    sql.NullTime    `db:"item_last_update"`
}

type CreateInvoiceRequest struct {
	UserId    uuid.UUID `json:"user_id"`
	Paid      bool      `json:"paid"`
	Total     float64   `json:"total"`
	CreatedBy string    `json:"created_by"`
}

type UpdateInvoiceRequest struct {
	Id            uuid.UUID `json:"id"`
	Paid          bool      `json:"paid"`
	Total         float64   `json:"total"`
	LastChangedBy string    `json:"last_changed_by"`
}

type ItemsToInvoiceRequest struct {
	InvoiceId uuid.UUID   `json:"invoice_id"`
	Items     []uuid.UUID `json:"items"`
}

type SimpleInvoiceItem struct {
	InvoiceId uuid.UUID `db:"invoice_id" json:"invoice_id"`
	ItemId    uuid.UUID `db:"item_id" json:"item_id"`
}

type ItemsToInvoiceResponse struct {
	InvoiceId uuid.UUID   `json:"invoice_id"`
	Items     []uuid.UUID `json:"items"`
	Success   bool        `json:"success"`
}

func (r *ItemsToInvoiceRequest) ToSimpleInvoiceItems() []SimpleInvoiceItem {
	var items []SimpleInvoiceItem
	for _, itemId := range r.Items {
		items = append(items, SimpleInvoiceItem{
			InvoiceId: r.InvoiceId,
			ItemId:    itemId,
		})
	}
	return items
}

type InvoiceRepository interface {
	CreateInvoice(request CreateInvoiceRequest) (InvoiceRow, error)
	UpdateInvoice(request UpdateInvoiceRequest) (InvoiceRow, error)
	DeleteInvoice(id uuid.UUID) (commons.DeleteResult, error)
	AddItemsToInvoice(request ItemsToInvoiceRequest) (ItemsToInvoiceResponse, error)
	RemoveItemFromInvoice(request SimpleInvoiceItem) (ItemsToInvoiceResponse, error)
	GetInvoice(id uuid.UUID) (InvoiceRow, error)
	GetInvoiceWithItems(id uuid.UUID) ([]InvoiceItemRow, error)
	GetAll(pagination *commons.Pagination) ([]InvoiceRow, error)
	GetAllForUser(userId uuid.UUID) ([]InvoiceRow, error)
}

type InvoiceRepositoryImpl struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{db: db}
}

const (
	CreateQuery                = `INSERT INTO invoices (user_id, total, paid, created_by) VALUES ($1, $2, $3, $4) RETURNING *`
	UpdateQuery                = `UPDATE invoices SET total = $2, paid = $3, last_changed_by = $4 WHERE alt_id = $1 RETURNING *`
	DeleteQuery                = `DELETE FROM invoices WHERE alt_id = $1`
	AddItemsToInvoiceQuery     = `INSERT INTO invoices_items (invoice_id, item_id) VALUES (:invoice_id, :item_id)`
	RemoveItemFromInvoiceQuery = `DELETE FROM invoices_items WHERE invoice_id = $1 AND item_id = $2`
	GetInvoiceQuery            = `SELECT * FROM invoices WHERE alt_id = $1`
	GetInvoiceWithItemsQuery   = `SELECT i.*, i2.id as item_seq, i2.alt_id as item_alt_id, i2.name as item_name, description as item_description, i2.unit_price as item_unit_price, i2.created_by as item_created_by, i2.created_at as item_created_at, i2.last_changed_by as item_last_changed_by, i2.last_update as item_last_update  FROM invoices i FULL OUTER JOIN invoices_items ii ON i.alt_id = ii.invoice_id FULL OUTER JOIN public.items i2 on i2.alt_id = ii.item_id WHERE i.alt_id = $1`
	GetAllQuery                = `SELECT * FROM invoices`
	GetAllWithPaginationQuery  = `SELECT * FROM invoices WHERE id > $1 LIMIT $2`
	GetAllForUserQuery         = `SELECT * FROM invoices WHERE user_id = $1`
)

func (r *InvoiceRepositoryImpl) CreateInvoice(request CreateInvoiceRequest) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, CreateQuery, request.UserId, request.Total, request.Paid, request.CreatedBy)
	return results, err
}

func (r *InvoiceRepositoryImpl) UpdateInvoice(request UpdateInvoiceRequest) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, UpdateQuery, request.Id, request.Total, request.Paid, request.LastChangedBy)
	return results, err
}

func (r *InvoiceRepositoryImpl) DeleteInvoice(id uuid.UUID) (commons.DeleteResult, error) {
	result, err := r.db.Exec(DeleteQuery, id)
	if err != nil {
		return commons.DeleteResult{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	return commons.DeleteResult{
		Id:      id,
		Deleted: rowsAffected > 0,
	}, nil
}

func (r *InvoiceRepositoryImpl) AddItemsToInvoice(request ItemsToInvoiceRequest) (ItemsToInvoiceResponse, error) {
	items := request.ToSimpleInvoiceItems()
	_, err := r.db.NamedExec(AddItemsToInvoiceQuery, items)
	if err != nil {
		return ItemsToInvoiceResponse{}, err
	} else {
		return ItemsToInvoiceResponse{
			InvoiceId: request.InvoiceId,
			Items:     request.Items,
			Success:   true,
		}, nil
	}
}

func (r *InvoiceRepositoryImpl) RemoveItemFromInvoice(request SimpleInvoiceItem) (ItemsToInvoiceResponse, error) {
	result, err := r.db.Exec(RemoveItemFromInvoiceQuery, request.InvoiceId, request.ItemId)
	if err != nil {
		return ItemsToInvoiceResponse{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	return ItemsToInvoiceResponse{
		InvoiceId: request.InvoiceId,
		Items:     []uuid.UUID{request.ItemId},
		Success:   rowsAffected > 0,
	}, nil
}

func (r *InvoiceRepositoryImpl) GetInvoice(id uuid.UUID) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, GetInvoiceQuery, id)
	return results, err
}

func (r *InvoiceRepositoryImpl) GetInvoiceWithItems(id uuid.UUID) ([]InvoiceItemRow, error) {
	var results []InvoiceItemRow
	err := r.db.Select(&results, GetInvoiceWithItemsQuery, id)
	return results, err
}

func (r *InvoiceRepositoryImpl) GetAll(pagination *commons.Pagination) ([]InvoiceRow, error) {
	var results []InvoiceRow
	var err error
	if pagination == nil {
		err = r.db.Select(&results, GetAllQuery)
	} else {
		err = r.db.Select(&results, GetAllWithPaginationQuery, pagination.LastId, pagination.PageSize)
	}
	return results, err
}

func (r *InvoiceRepositoryImpl) GetAllForUser(userId uuid.UUID) ([]InvoiceRow, error) {
	var results []InvoiceRow
	err := r.db.Select(&results, GetAllForUserQuery, userId)
	return results, err
}
