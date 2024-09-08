package invoice

import (
	"github.com/google/uuid"
	"inventory-service-go/commons"
	"inventory-service-go/item"
	"time"
)

type Invoice struct {
	Seq       int               `json:"seq"`
	Id        uuid.UUID         `json:"id"`
	UserId    uuid.UUID         `json:"user_id"`
	Total     float64           `json:"total"`
	Paid      bool              `json:"paid"`
	Items     []item.Item       `json:"items"`
	AuditInfo commons.AuditInfo `json:"audit_info"`
}

func fromRow(row InvoiceRow) Invoice {
	return Invoice{
		Seq:    int(row.Id),
		Id:     row.AltId,
		UserId: row.UserId,
		Total:  row.Total,
		Paid:   row.Paid,
		Items:  []item.Item{},
		AuditInfo: commons.AuditInfo{
			CreatedBy:     row.CreatedBy,
			CreatedAt:     row.CreatedAt.Format(time.RFC3339),
			LastChangedBy: row.LastChangedBy,
			LastUpdate:    row.LastUpdate.Format(time.RFC3339),
		},
	}
}

func fromRows(results []InvoiceRow) []Invoice {
	invoices := []Invoice{}
	for _, row := range results {
		invoices = append(invoices, fromRow(row))
	}
	return invoices
}

func fromRowWithItems(row []InvoiceItemRow) Invoice {
	var items []item.Item
	for _, row := range row {
		if row.ItemSeqId.Valid == false {
			continue
		}
		items = append(items, item.Item{
			Seq:         int(row.ItemSeqId.Int64),
			Id:          row.ItemAltId,
			Name:        row.ItemName.String,
			Description: row.ItemDescription.String,
			UnitPrice:   row.ItemUnitPrice.Float64,
			AuditInfo: commons.AuditInfo{
				CreatedBy:     row.ItemCreatedBy.String,
				CreatedAt:     row.ItemCreatedAt.Time.Format(time.RFC3339),
				LastUpdate:    row.ItemLastUpdate.Time.Format(time.RFC3339),
				LastChangedBy: row.ItemLastChangedBy.String,
			},
		})
	}
	return Invoice{
		Seq:    int(row[0].Id),
		Id:     row[0].AltId,
		UserId: row[0].UserId,
		Total:  row[0].Total,
		Paid:   row[0].Paid,
		Items:  items,
		AuditInfo: commons.AuditInfo{
			CreatedBy:     row[0].CreatedBy,
			CreatedAt:     row[0].CreatedAt.Format(time.RFC3339),
			LastChangedBy: row[0].LastChangedBy,
			LastUpdate:    row[0].LastUpdate.Format(time.RFC3339),
		},
	}
}

type InvoiceService interface {
	GetInvoice(id uuid.UUID, withItems bool) (Invoice, error)
	GetInvoicesForUser(userId uuid.UUID) ([]Invoice, error)
	CreateInvoice(invoice CreateInvoiceRequest) (Invoice, error)
	UpdateInvoice(invoice UpdateInvoiceRequest) (Invoice, error)
	DeleteInvoice(id uuid.UUID) (commons.DeleteResult, error)
	GetAllInvoices(pagination *commons.Pagination) ([]Invoice, error)
	AddItemsToInvoice(request ItemsToInvoiceRequest) (ItemsToInvoiceResponse, error)
	RemoveItemFromInvoice(request SimpleInvoiceItem) (ItemsToInvoiceResponse, error)
}

type InvoiceServiceImpl struct {
	repo InvoiceRepository
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceServiceImpl {
	return &InvoiceServiceImpl{
		repo: repo,
	}
}

func (s *InvoiceServiceImpl) GetInvoice(id uuid.UUID, withItems bool) (Invoice, error) {
	if withItems {
		results, err := s.repo.GetInvoiceWithItems(id)
		invoice := fromRowWithItems(results)
		return invoice, err
	} else {
		results, err := s.repo.GetInvoice(id)
		invoice := fromRow(results)
		return invoice, err
	}
}

func (s *InvoiceServiceImpl) GetInvoicesForUser(userId uuid.UUID) ([]Invoice, error) {
	invoices, err := s.repo.GetAllForUser(userId)
	if err != nil {
		return nil, err
	}
	result := fromRows(invoices)
	return result, nil
}

func (s *InvoiceServiceImpl) CreateInvoice(invoice CreateInvoiceRequest) (Invoice, error) {
	invoiceRow, err := s.repo.CreateInvoice(invoice)
	if err != nil {
		return Invoice{}, err
	}
	return fromRow(invoiceRow), nil
}

func (s *InvoiceServiceImpl) UpdateInvoice(invoice UpdateInvoiceRequest) (Invoice, error) {
	invoiceRow, err := s.repo.UpdateInvoice(invoice)
	if err != nil {
		return Invoice{}, err
	}
	return fromRow(invoiceRow), nil
}

func (s *InvoiceServiceImpl) DeleteInvoice(id uuid.UUID) (commons.DeleteResult, error) {
	results, err := s.repo.DeleteInvoice(id)
	if err != nil {
		return commons.DeleteResult{}, err
	}
	return results, nil
}

func (s *InvoiceServiceImpl) GetAllInvoices(pagination *commons.Pagination) ([]Invoice, error) {
	results, err := s.repo.GetAll(pagination)
	if err != nil {
		return nil, err
	}
	invoices := fromRows(results)
	return invoices, nil
}

func (s *InvoiceServiceImpl) AddItemsToInvoice(request ItemsToInvoiceRequest) (ItemsToInvoiceResponse, error) {
	results, err := s.repo.AddItemsToInvoice(request)
	if err != nil {
		return ItemsToInvoiceResponse{}, err
	}
	return results, err
}

func (s *InvoiceServiceImpl) RemoveItemFromInvoice(request SimpleInvoiceItem) (ItemsToInvoiceResponse, error) {
	results, err := s.repo.RemoveItemFromInvoice(request)
	if err != nil {
		return ItemsToInvoiceResponse{}, err
	}
	return results, err
}
