package invoice

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
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
