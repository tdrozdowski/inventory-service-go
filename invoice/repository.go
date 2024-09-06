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
	Paid          bool    `json:"paid"`
	Total         float64 `json:"total"`
	LastChangedBy string  `json:"last_changed_by"`
}

type AddItemsToInvoiceRequest struct {
	InvoiceId uuid.UUID   `json:"invoice_id"`
	Items     []uuid.UUID `json:"items"`
}

type InvoiceRepository interface {
	CreateInvoice(request CreateInvoiceRequest) (InvoiceRow, error)
	UpdateInvoice(request UpdateInvoiceRequest) (InvoiceRow, error)
	AddItemsToInvoice(request AddItemsToInvoiceRequest) ([]InvoiceItemRow, error)
	GetInvoice(id int64) (InvoiceRow, error)
	GetInvoiceItems(id uuid.UUID) ([]InvoiceItemRow, error)
	GetInvoiceItem(id uuid.UUID) (InvoiceItemRow, error)
	GetInvoices(id uuid.UUID) ([]InvoiceRow, error)
	GetInvoicesByUserId(userId uuid.UUID) ([]InvoiceRow, error)
}

type InvoiceRepositoryImpl struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) *InvoiceRepositoryImpl {
	return &InvoiceRepositoryImpl{db: db}
}

const (
	CREATE_QUERY = `INSERT INTO invoices (user_id, total, paid, created_by) VALUES ($1, $2, $3, $4) RETURNING *`
)

func (r *InvoiceRepositoryImpl) CreateInvoice(request CreateInvoiceRequest) (InvoiceRow, error) {
	var results = InvoiceRow{}
	err := r.db.Get(&results, CREATE_QUERY, request.UserId, request.Paid, request.Total, request.CreatedBy)
	return results, err
}
