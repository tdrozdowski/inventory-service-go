package person

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"inventory-service-go/commons"
	"testing"
)

func TestPersonRepositoryImpl_GetAll(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		pagination *commons.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success with pagination",
			fields: fields{
				db: nil,
			},
			args: args{
				pagination: &commons.Pagination{
					PageSize: 10,
					LastId:   1,
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("SELECT \\* FROM persons WHERE id > \\? LIMIT \\?").
					WithArgs(1, 10).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Success without pagination",
			fields: fields{
				db: nil,
			},
			args: args{
				pagination: nil,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("SELECT \\* FROM persons").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "DB error with pagination",
			fields: fields{
				db: nil,
			},
			args: args{
				pagination: &commons.Pagination{
					PageSize: 10,
					LastId:   1,
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM persons WHERE id > \\? LIMIT \\?").
					WithArgs(1, 10).WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
		{
			name: "DB error without pagination",
			fields: fields{
				db: nil,
			},
			args: args{
				pagination: nil,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM persons").
					WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			_, err := p.GetAll(tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepositoryImpl.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
