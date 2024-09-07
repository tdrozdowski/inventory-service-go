package invoice

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"inventory-service-go/commons"
	"testing"
	"time"
)

func TestInvoiceRepositoryImpl_CreateInvoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	now := time.Now()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	newUuid := uuid.New()
	testCases := []struct {
		name    string
		request CreateInvoiceRequest
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "Successful Invoice Creation",
			request: CreateInvoiceRequest{
				UserId:    uuid.New(),
				Paid:      true,
				Total:     123.45,
				CreatedBy: "test_user",
			},
			rows: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(1, newUuid, uuid.New(), true, 123.45, "test_user", now, now, "test_user"),
			wantErr: false,
		},
		{
			name: "Failed Invoice Creation",
			request: CreateInvoiceRequest{
				UserId:    uuid.New(),
				Paid:      false,
				Total:     0.0,
				CreatedBy: "test_user",
			},
			rows:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.rows == nil {
				mock.ExpectQuery("INSERT INTO invoices").
					WithArgs(tc.request.UserId, tc.request.Paid, tc.request.Total, tc.request.CreatedBy).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("INSERT INTO invoices").
					WithArgs(tc.request.UserId, tc.request.Paid, tc.request.Total, tc.request.CreatedBy).
					WillReturnRows(tc.rows)
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			results, err := r.CreateInvoice(tc.request)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, results)
				assert.Equal(t, results.Id, int64(1))
				assert.Equal(t, results.AltId, newUuid)
			}
		})
	}
}

func TestInvoiceRepositoryImpl_UpdateInvoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	now := time.Now()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	newUuid := uuid.New()
	testCases := []struct {
		name    string
		request UpdateInvoiceRequest
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "Successful Invoice Update",
			request: UpdateInvoiceRequest{
				Id:            newUuid,
				Paid:          true,
				Total:         123.45,
				LastChangedBy: "updated_user",
			},
			rows: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(1, newUuid, uuid.New(), true, 123.45, "created_user", now, now, "updated_user"),
			wantErr: false,
		},
		{
			name: "Failed Invoice Update",
			request: UpdateInvoiceRequest{
				Id:            newUuid,
				Paid:          false,
				Total:         0.0,
				LastChangedBy: "update_failed_user",
			},
			rows:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.rows == nil {
				mock.ExpectQuery("UPDATE invoices").
					WithArgs(tc.request.Id, tc.request.Total, tc.request.Paid, tc.request.LastChangedBy).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("UPDATE invoices").
					WithArgs(tc.request.Id, tc.request.Total, tc.request.Paid, tc.request.LastChangedBy).
					WillReturnRows(tc.rows)
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			results, err := r.UpdateInvoice(tc.request)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, results)
				assert.Equal(t, results.Id, int64(1))
				assert.Equal(t, results.AltId, newUuid)
			}
		})
	}
}

func TestInvoiceRepositoryImpl_AddItemsToInvoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	invoiceId := uuid.New()
	itemId1 := uuid.New()
	itemId2 := uuid.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	testCases := []struct {
		name    string
		request ItemsToInvoiceRequest
		items   []SimpleInvoiceItem
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "Successful Items Addition",
			request: ItemsToInvoiceRequest{
				InvoiceId: invoiceId,
				Items:     []uuid.UUID{itemId1, itemId2},
			},
			items: []SimpleInvoiceItem{
				{
					InvoiceId: invoiceId,
					ItemId:    itemId1,
				},
				{
					InvoiceId: invoiceId,
					ItemId:    itemId2,
				},
			},
			rows: sqlmock.NewRows([]string{"invoice_id", "item_id"}).
				AddRow(invoiceId, itemId1).
				AddRow(invoiceId, itemId2),
			wantErr: false,
		},
		{
			name: "Failed Items Addition",
			request: ItemsToInvoiceRequest{
				InvoiceId: uuid.New(),
				Items:     []uuid.UUID{uuid.New(), uuid.New()},
			},
			items: []SimpleInvoiceItem{
				{
					InvoiceId: invoiceId,
					ItemId:    itemId1,
				},
				{
					InvoiceId: invoiceId,
					ItemId:    itemId2,
				},
			},
			rows:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr {
				mock.ExpectExec("INSERT INTO invoices_items").
					WithArgs(tc.items).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectExec("INSERT INTO invoices_items").
					WithArgs(tc.items[0].InvoiceId, tc.items[0].ItemId, tc.items[1].InvoiceId, tc.items[1].ItemId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			results, err := r.AddItemsToInvoice(tc.request)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, results)
				assert.Equal(t, results.InvoiceId, tc.request.InvoiceId)
				assert.Equal(t, results.Items, tc.request.Items)
				assert.True(t, results.Success)
			}
		})
	}
}
func TestInvoiceRepositoryImpl_GetInvoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	invoiceId := uuid.New()

	testCases := []struct {
		name    string
		id      uuid.UUID
		row     *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "Successful Invoice Fetching",
			id:   invoiceId,
			row: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(1, invoiceId, uuid.New(), true, 123.45, "test_user", time.Now(), time.Now(), "test_user"),
			wantErr: false,
		},
		{
			name:    "Failed Invoice Fetching",
			id:      uuid.New(),
			row:     nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.row == nil {
				mock.ExpectQuery("SELECT .* FROM invoices WHERE alt_id").
					WithArgs(tc.id).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("SELECT .* FROM invoices WHERE alt_id").
					WithArgs(tc.id).
					WillReturnRows(tc.row)
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			result, err := r.GetInvoice(tc.id)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, result.AltId, invoiceId)
			}
		})
	}
}

func TestInvoiceRepositoryImpl_GetInvoiceWithItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	invoiceId := uuid.New()
	userId := uuid.New()
	now := time.Now()
	testCases := []struct {
		name    string
		id      uuid.UUID
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "Successful Getting invoice with items",
			id:   invoiceId,
			rows: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "total", "paid", "created_by", "created_at", "last_changed_by", "last_update", "item_name", "item_description", "item_unit_price", "item_created_by", "item_created_at", "item_last_changed_by", "item_last_update"}).
				AddRow(1, invoiceId, userId, 100.0, true, "unit_test", now, "unit_test", now, "Item1", "Item 1", 12.34, "unit_test", now, "unit_test", now).
				AddRow(1, invoiceId, userId, 100.0, true, "unit_test", now, "unit_test", now, "Item2", "Item 2", 56.78, "unit_test", now, "unit_test", now),
			wantErr: false,
		},
		{
			name:    "Failed Getting invoice with items",
			id:      uuid.New(),
			rows:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.rows == nil {
				mock.ExpectQuery("SELECT .* FROM invoices_items WHERE invoice_id").
					WithArgs(tc.id).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("SELECT i.*, i2.alt_id as item_alt_id, i2.name as item_name, description as item_description, i2.unit_price as item_unit_price, i2.created_by as item_created_by, i2.created_at as item_created_at, i2.last_changed_by as item_last_changed_by, i2.last_update as item_last_update FROM invoices i INNER JOIN invoices_items ii ON i.id = ii.invoice_id INNER JOIN public.items i2 on i2.alt_id = ii.item_id WHERE i.alt_id = \\$1").
					WithArgs(tc.id).
					WillReturnRows(tc.rows)
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			results, err := r.GetInvoiceWithItems(tc.id)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, results)
				assert.Equal(t, len(results), 2)
				assert.Equal(t, results[0].AltId, invoiceId)
				assert.Equal(t, results[0].ItemName, "Item1")
				assert.Equal(t, results[1].ItemName, "Item2")
			}
		})
	}
}

func TestInvoiceRepositoryImpl_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	testCases := []struct {
		name           string
		pagination     *commons.Pagination
		rows           *sqlmock.Rows
		expectedLength int
		wantErr        bool
	}{
		{
			name:       "Successful Fetching All Invoices Without Pagination",
			pagination: nil,
			rows: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(1, uuid.New(), uuid.New(), true, 123.45, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(2, uuid.New(), uuid.New(), false, 543.21, "test_user", time.Now(), time.Now(), "test_user"),
			expectedLength: 2,
			wantErr:        false,
		},
		{
			name:       "Successful Fetching All Invoices With Pagination",
			pagination: &commons.Pagination{LastId: 1, PageSize: 5},
			rows: sqlmock.NewRows([]string{"id", "alt_id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(2, uuid.New(), uuid.New(), true, 123.45, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(3, uuid.New(), uuid.New(), false, 543.21, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(4, uuid.New(), uuid.New(), true, 123.45, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(5, uuid.New(), uuid.New(), false, 543.21, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(6, uuid.New(), uuid.New(), true, 123.45, "test_user", time.Now(), time.Now(), "test_user"),
			expectedLength: 5,
			wantErr:        false,
		},
		{
			name:       "Failed Fetching All Invoices",
			pagination: &commons.Pagination{LastId: 1, PageSize: 5},
			rows:       nil,
			wantErr:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.rows == nil {
				mock.ExpectQuery("SELECT .* FROM invoices").
					WithArgs().
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("SELECT .* FROM invoices").
					WithArgs().
					WillReturnRows(tc.rows)
			}
			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))
			result, err := r.GetAll(tc.pagination)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(result), tc.expectedLength)
			}
		})
	}
}

func TestInvoiceRepositoryImpl_GetAllForUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	userId := uuid.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	testCases := []struct {
		name           string
		id             uuid.UUID
		rows           *sqlmock.Rows
		expectedLength int
		wantErr        bool
	}{
		{
			name: "Successful Fetching All User Invoices",
			id:   uuid.New(),
			rows: sqlmock.NewRows([]string{"id", "user_id", "paid", "total", "created_by", "created_at", "last_update", "last_changed_by"}).
				AddRow(1, userId, true, 123.45, "test_user", time.Now(), time.Now(), "test_user").
				AddRow(2, userId, false, 543.21, "test_user", time.Now(), time.Now(), "test_user"),
			expectedLength: 2,
			wantErr:        false,
		},
		{
			name:    "Failed Fetching All User Invoices",
			id:      uuid.New(),
			rows:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr && tc.rows == nil {
				mock.ExpectQuery("SELECT * FROM invoices WHERE user_id = \\$1").
					WithArgs(userId).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectQuery("SELECT * FROM invoices WHERE user_id = $1").
					WithArgs(userId).
					WillReturnRows(tc.rows)
			}
			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))
			result, err := r.GetAllForUser(userId)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(result), tc.expectedLength)
			}
		})
	}
}
func TestInvoiceRepositoryImpl_RemoveItemFromInvoice(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	invoiceId := uuid.New()
	itemId := uuid.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	testCases := []struct {
		name    string
		request SimpleInvoiceItem
		rows    driver.Result
		wantErr bool
	}{
		{
			name: "Successful Removal of Item",
			request: SimpleInvoiceItem{
				InvoiceId: invoiceId,
				ItemId:    itemId,
			},
			rows:    sqlmock.NewResult(1, 1),
			wantErr: false,
		},
		{
			name: "Failed Item Removal",
			request: SimpleInvoiceItem{
				InvoiceId: invoiceId,
				ItemId:    itemId,
			},
			rows:    sqlmock.NewErrorResult(errors.New("error")),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr {
				mock.ExpectExec("DELETE FROM invoices_items WHERE invoice_id = $1 AND item_id = $2").
					WithArgs(tc.request.InvoiceId, tc.request.ItemId).
					WillReturnError(errors.New("error"))
			} else {
				mock.ExpectExec("DELETE FROM invoices_items WHERE invoice_id = $1 AND item_id = $2").
					WithArgs(tc.request.InvoiceId, tc.request.ItemId).
					WillReturnResult(tc.rows)
			}

			r := NewInvoiceRepository(sqlx.NewDb(db, "mockDb"))

			results, err := r.RemoveItemFromInvoice(tc.request)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.True(t, results.Success)
				assert.Equal(t, results.InvoiceId, tc.request.InvoiceId)
				assert.Equal(t, results.Items[0], tc.request.ItemId)
			}
		})
	}
}
