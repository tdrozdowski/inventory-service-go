package invoice

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	Id                int64     `db:"id"`
	AltId             uuid.UUID `db:"alt_id"`
	UserId            uuid.UUID `db:"user_id"`
	Total             float64   `db:"total"`
	Paid              bool      `db:"paid"`
	CreatedBy         string    `db:"created_by"`
	CreatedAt         time.Time `db:"created_at"`
	LastChangedBy     string    `db:"last_changed_by"`
	LastUpdate        time.Time `db:"last_update"`
	ItemAltId         uuid.UUID `db:"item_alt_id"`
	ItemName          string    `db:"item_name"`
	ItemDescription   string    `db:"item_description"`
	ItemUnitPrice     float64   `db:"item_unit_price"`
	ItemCreatedBy     string    `db:"item_created_by"`
	ItemCreatedAt     time.Time `db:"item_created_at"`
	ItemLastChangedBy string    `db:"item_last_changed_by"`
	ItemLastUpdate    time.Time `db:"item_last_update"`
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

type AddItemsToInvoiceRequest struct {
	InvoiceId uuid.UUID   `json:"invoice_id"`
	Items     []uuid.UUID `json:"items"`
}

type SimpleInvoiceItem struct {
	InvoiceId uuid.UUID `db:"invoice_id"`
	ItemId    uuid.UUID `db:"item_id"`
}

type AddItemsToInvoiceResponse struct {
	InvoiceId uuid.UUID   `json:"invoice_id"`
	Items     []uuid.UUID `json:"items"`
	Success   bool        `json:"success"`
}

func (r *AddItemsToInvoiceRequest) ToSimpleInvoiceItems() []SimpleInvoiceItem {
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
	AddItemsToInvoice(request AddItemsToInvoiceRequest) (AddItemsToInvoiceResponse, error)
	GetInvoice(id uuid.UUID) (InvoiceRow, error)
	GetInvoiceWithItems(id uuid.UUID) ([]InvoiceItemRow, error)
	GetAll(id uuid.UUID) ([]InvoiceRow, error)
	GetAllForUser(userId uuid.UUID) ([]InvoiceRow, error)
}

type InvoiceRepositoryImpl struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{db: db}
}

const (
	CreateQuery              = `INSERT INTO invoices (user_id, total, paid, created_by) VALUES ($1, $2, $3, $4) RETURNING *`
	UpdateQuery              = `UPDATE invoices SET total = $2, paid = $3, last_changed_by = $4 WHERE alt_id = $1 RETURNING *`
	AddItemsToInvoiceQuery   = `INSERT INTO invoices_items (invoice_id, item_id) VALUES (:invoice_id, :item_id)`
	GetInvoiceQuery          = `SELECT * FROM invoices WHERE alt_id = $1`
	GetInvoiceWithItemsQuery = `SELECT i.*, i2.alt_id as item_alt_id, i2.name as item_name, description as item_description, i2.unit_price as item_unit_price, i2.created_by as item_created_by, i2.created_at as item_created_at, i2.last_changed_by as item_last_changed_by, i2.last_update as item_last_update  FROM invoices i INNER JOIN invoices_items ii ON i.id = ii.invoice_id INNER JOIN public.items i2 on i2.alt_id = ii.item_id WHERE i.alt_id = $1`
)

func (r *InvoiceRepositoryImpl) CreateInvoice(request CreateInvoiceRequest) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, CreateQuery, request.UserId, request.Paid, request.Total, request.CreatedBy)
	return results, err
}

func (r *InvoiceRepositoryImpl) UpdateInvoice(request UpdateInvoiceRequest) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, UpdateQuery, request.Id, request.Total, request.Paid, request.LastChangedBy)
	return results, err
}

func (r *InvoiceRepositoryImpl) AddItemsToInvoice(request AddItemsToInvoiceRequest) (AddItemsToInvoiceResponse, error) {
	items := request.ToSimpleInvoiceItems()
	_, err := r.db.NamedExec(AddItemsToInvoiceQuery, items)
	if err != nil {
		return AddItemsToInvoiceResponse{}, err
	} else {
		return AddItemsToInvoiceResponse{
			InvoiceId: request.InvoiceId,
			Items:     request.Items,
			Success:   true,
		}, nil
	}
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
