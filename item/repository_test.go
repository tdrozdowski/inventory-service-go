package item

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"inventory-service-go/commons"
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

func TestItemRepository_UpdateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	itemtest := ItemRow{
		AltId:       uuid.New(),
		Name:        "TestItem",
		Description: "This is a test item",
		UnitPrice:   20.0,
		CreatedBy:   "testUser",
	}

	itemtestUpd := ItemRow{
		AltId:         itemtest.AltId,
		Name:          "UpdatedTestItem",
		Description:   "This is an updated test item",
		UnitPrice:     22.0,
		LastChangedBy: "testUser2",
	}

	updateQuery := "UPDATE items SET name = \\$1, description = \\$2, unit_price = \\$3, last_changed_by = \\$4 WHERE alt_id = \\$5 returning *"

	mock.ExpectQuery(updateQuery).
		WithArgs(itemtestUpd.Name, itemtestUpd.Description, itemtestUpd.UnitPrice, itemtestUpd.LastChangedBy, itemtest.AltId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "alt_id", "name", "description", "unit_price", "created_by", "created_at", "last_changed_by", "last_update"}).
				AddRow(1, itemtest.AltId, itemtestUpd.Name, itemtestUpd.Description, itemtestUpd.UnitPrice, itemtest.CreatedBy, time.Now(), itemtestUpd.LastChangedBy, time.Now()))

	itemRepo := NewItemRepository(sqlx.NewDb(db, ""))

	request := UpdateItemRequest{
		Id:            itemtest.AltId,
		Name:          itemtestUpd.Name,
		Description:   itemtestUpd.Description,
		UnitPrice:     itemtestUpd.UnitPrice,
		LastChangedBy: itemtestUpd.LastChangedBy,
	}

	row, err := itemRepo.UpdateItem(request)
	if err != nil {
		t.Errorf("error was not expected when updating item: %s", err)
	} else {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	}
	assert.Equal(t, itemtestUpd.Name, row.Name)
	assert.Equal(t, itemtestUpd.Description, row.Description)
	assert.Equal(t, itemtestUpd.UnitPrice, row.UnitPrice)
	assert.Equal(t, itemtestUpd.LastChangedBy, row.LastChangedBy)
	assert.Equal(t, itemtest.AltId, row.AltId)
	assert.Equal(t, itemtest.CreatedBy, row.CreatedBy)
	assert.Equal(t, itemtestUpd.LastChangedBy, row.LastChangedBy)
}

func TestItemRepository_GetItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	itemtest := ItemRow{
		AltId:       uuid.New(),
		Name:        "GetTestItem",
		Description: "This is a test item for get",
		UnitPrice:   30.0,
		CreatedBy:   "getTestUser",
	}

	rows := sqlmock.NewRows([]string{"id", "alt_id", "name", "description", "unit_price", "created_by", "created_at", "last_changed_by", "last_update"}).
		AddRow(1, itemtest.AltId, itemtest.Name, itemtest.Description, itemtest.UnitPrice, itemtest.CreatedBy, time.Now(), itemtest.CreatedBy, time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM items WHERE alt_id = \\$1$").
		WithArgs(itemtest.AltId).
		WillReturnRows(rows)

	itemRepo := NewItemRepository(sqlx.NewDb(db, ""))

	resultItem, err := itemRepo.GetItem(itemtest.AltId)
	if err != nil {
		t.Errorf("error was not expected when getting item: %s", err)
	} else {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	}

	assert.Equal(t, itemtest.Name, resultItem.Name)
	assert.Equal(t, itemtest.AltId, resultItem.AltId)
	assert.Equal(t, itemtest.Description, resultItem.Description)
	assert.Equal(t, itemtest.UnitPrice, resultItem.UnitPrice)
	assert.Equal(t, itemtest.CreatedBy, resultItem.CreatedBy)
}

func TestItemRepositoryImpl_GetItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	itemtest1 := ItemRow{
		AltId:       uuid.New(),
		Name:        "TestItem1",
		Description: "This is a test item 1",
		UnitPrice:   10.0,
		CreatedBy:   "testUser",
	}

	itemtest2 := ItemRow{
		AltId:       uuid.New(),
		Name:        "TestItem2",
		Description: "This is a test item 2",
		UnitPrice:   15.0,
		CreatedBy:   "testUser",
	}

	rows := sqlmock.NewRows([]string{"alt_id", "name", "description", "unit_price", "created_by", "created_at", "last_changed_by", "last_update"}).
		AddRow(itemtest1.AltId, itemtest1.Name, itemtest1.Description, itemtest1.UnitPrice, itemtest1.CreatedBy, time.Now(), itemtest1.CreatedBy, time.Now()).
		AddRow(itemtest2.AltId, itemtest2.Name, itemtest2.Description, itemtest2.UnitPrice, itemtest2.CreatedBy, time.Now(), itemtest2.CreatedBy, time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM items$").
		WillReturnRows(rows)

	itemRepo := NewItemRepository(sqlx.NewDb(db, ""))
	items, err := itemRepo.GetItems(nil)
	if err != nil {
		t.Errorf("error was not expected when getting items: %s", err)
	} else {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	}

	assert.Equal(t, 2, len(items))
	assert.Equal(t, itemtest1.AltId, items[0].AltId)
	assert.Equal(t, itemtest1.Name, items[0].Name)
	assert.Equal(t, itemtest1.Description, items[0].Description)
	assert.Equal(t, itemtest1.UnitPrice, items[0].UnitPrice)
	assert.Equal(t, itemtest1.CreatedBy, items[0].CreatedBy)

	assert.Equal(t, itemtest2.AltId, items[1].AltId)
	assert.Equal(t, itemtest2.Name, items[1].Name)
	assert.Equal(t, itemtest2.Description, items[1].Description)
	assert.Equal(t, itemtest2.UnitPrice, items[1].UnitPrice)
	assert.Equal(t, itemtest2.CreatedBy, items[1].CreatedBy)
}

func TestItemRepositoryImpl_GetAllWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	itemtest1 := ItemRow{
		AltId:       uuid.New(),
		Name:        "TestItem1",
		Description: "This is a test item 1",
		UnitPrice:   10.0,
		CreatedBy:   "testUser",
	}

	itemtest2 := ItemRow{
		AltId:       uuid.New(),
		Name:        "TestItem2",
		Description: "This is a test item 2",
		UnitPrice:   15.0,
		CreatedBy:   "testUser",
	}

	rows := sqlmock.NewRows([]string{"alt_id", "name", "description", "unit_price", "created_by", "created_at", "last_changed_by", "last_update"}).
		AddRow(itemtest2.AltId, itemtest2.Name, itemtest2.Description, itemtest2.UnitPrice, itemtest2.CreatedBy, time.Now(), itemtest2.CreatedBy, time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM items WHERE id > \\$1 LIMIT \\$2$").
		WillReturnRows(rows)

	itemRepo := NewItemRepository(sqlx.NewDb(db, ""))
	pagination := commons.Pagination{
		LastId:   int(itemtest1.Id),
		PageSize: 1,
	}
	items, err := itemRepo.GetItems(&pagination)
	if err != nil {
		t.Errorf("error was not expected when getting items: %s", err)
	} else {
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	}

	assert.Equal(t, 1, len(items))

	assert.Equal(t, itemtest2.AltId, items[0].AltId)
	assert.Equal(t, itemtest2.Name, items[0].Name)
	assert.Equal(t, itemtest2.Description, items[0].Description)
	assert.Equal(t, itemtest2.UnitPrice, items[0].UnitPrice)
	assert.Equal(t, itemtest2.CreatedBy, items[0].CreatedBy)
}

func TestItemRepositoryImpl_DeleteItem(t *testing.T) {
	var testCases = []struct {
		testName string
		id       uuid.UUID
		wantErr  bool
	}{
		{"Test Delete Item - Success", uuid.New(), false},
		{"Test Delete Item - Fail", uuid.Nil, true},
	}
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			itemRepo := NewItemRepository(sqlx.NewDb(db, ""))

			if !tt.wantErr {
				mock.ExpectQuery("^DELETE FROM items WHERE alt_id = \\$1").
					WithArgs(tt.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "deleted"}).AddRow(tt.id, true))
			} else {
				mock.ExpectQuery("^DELETE FROM items WHERE alt_id = \\$1").
					WithArgs(tt.id).
					WillReturnError(sql.ErrNoRows)
			}

			_, err = itemRepo.DeleteItem(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
