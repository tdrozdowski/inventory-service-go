package item

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestItemRepository_CreateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	itemtest := ItemRow{
		Name:        "TestItem",
		Description: "This is a test item",
		UnitPrice:   10.0,
		CreatedBy:   "testUser",
	}

	newUuid := uuid.NewString()

	rows := sqlmock.NewRows([]string{"id", "alt_id", "name", "description", "unit_price", "created_by", "created_at", "last_changed_by", "last_update"}).
		AddRow(1, newUuid, itemtest.Name, itemtest.Description, itemtest.UnitPrice, itemtest.CreatedBy, time.Now(), itemtest.CreatedBy, time.Now())

	mock.ExpectQuery("^INSERT INTO items (.+) VALUES (.+)$").
		WithArgs(itemtest.Name, itemtest.Description, itemtest.UnitPrice, itemtest.CreatedBy).
		WillReturnRows(rows)

	itemRepo := NewItemRepository(sqlx.NewDb(db, ""))
	request := CreateItemRequest{Name: itemtest.Name, Description: itemtest.Description, UnitPrice: itemtest.UnitPrice, CreatedBy: itemtest.CreatedBy}
	resultItem, err := itemRepo.CreateItem(request)
	if err != nil {
		t.Errorf("error was not expected when creating item: %s", err)
	}

	assert.Equal(t, itemtest.Name, resultItem.Name)
	assert.Equal(t, newUuid, resultItem.AltId.String())
	assert.Equal(t, itemtest.Description, resultItem.Description)
	assert.Equal(t, itemtest.UnitPrice, resultItem.UnitPrice)
	assert.Equal(t, itemtest.CreatedBy, resultItem.CreatedBy)
	assert.Equal(t, itemtest.CreatedBy, resultItem.LastChangedBy)
}
