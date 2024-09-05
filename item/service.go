package item

import (
	"github.com/google/uuid"
	"inventory-service-go/commons"
)

type Item struct {
	Seq         int               `json:"seq"`
	Id          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	UnitPrice   float64           `json:"unit_price"`
	AuditInfo   commons.AuditInfo `json:"audit_info"`
}

func itemFromRow(row ItemRow) Item {
	return Item{
		Seq:         int(row.Id),
		Id:          row.AltId,
		Name:        row.Name,
		Description: row.Description,
		UnitPrice:   row.UnitPrice,
		AuditInfo: commons.AuditInfo{
			CreatedBy:     row.CreatedBy,
			CreatedAt:     row.CreatedAt,
			LastUpdate:    row.LastUpdate,
			LastChangedBy: row.LastChangedBy,
		},
	}
}

type ItemService interface {
	CreateItem(request CreateItemRequest) (*Item, error)
	UpdateItem(request UpdateItemRequest) (*Item, error)
	DeleteItem(id uuid.UUID) (*Item, error)
	GetItem(id uuid.UUID) (*Item, error)
	GetItems(pagination commons.Pagination) ([]Item, error)
}

type ItemServiceImpl struct {
	repo ItemRepository
}

func NewItemService(repo ItemRepository) *ItemServiceImpl {
	return &ItemServiceImpl{
		repo: repo,
	}
}

func (s *ItemServiceImpl) CreateItem(request CreateItemRequest) (*Item, error) {
	row, err := s.repo.CreateItem(request)
	if err != nil {
		return nil, err
	}
	i := itemFromRow(row)
	return &i, nil
}
