package person

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
				mock.ExpectQuery("SELECT \\* FROM persons WHERE id > \\$1 LIMIT \\$2").
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
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)
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

func TestPersonRepositoryImpl_GetById(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			fields: fields{
				db: nil,
			},
			args: args{
				id: 1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("SELECT \\* FROM persons WHERE id = \\$1").
					WithArgs(1).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "DB error",
			fields: fields{
				db: nil,
			},
			args: args{
				id: 1,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM persons WHERE id = \\$1").
					WithArgs(1).WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			_, err := p.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepositoryImpl.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPersonRepositoryImpl_GetByUuid(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		uuid uuid.UUID
	}
	testUuid, _ := uuid.Parse("2b1b425e-dee2-4227-8d94-f470a0ce0cd0")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			fields: fields{
				db: nil,
			},
			args: args{
				uuid: testUuid,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("SELECT \\* FROM persons WHERE alt_id = \\$1").
					WithArgs("2b1b425e-dee2-4227-8d94-f470a0ce0cd0").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "DB error",
			fields: fields{
				db: nil,
			},
			args: args{
				uuid: testUuid,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM persons WHERE alt_id = \\$1").
					WithArgs("2b1b425e-dee2-4227-8d94-f470a0ce0cd0").WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			_, err := p.GetByUuid(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepositoryImpl.GetByUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPersonRepositoryImpl_Create(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		request CreatePersonRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			fields: fields{
				db: nil,
			},
			args: args{
				request: CreatePersonRequest{
					Name:      "test name",
					Email:     "test email",
					CreatedBy: "test user",
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("INSERT INTO persons ").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "DB error",
			fields: fields{
				db: nil,
			},
			args: args{
				request: CreatePersonRequest{
					Name:      "test name",
					Email:     "test email",
					CreatedBy: "test user",
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO persons ").
					WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			_, err := p.Create(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepositoryImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPersonRepositoryImpl_DeleteByUuid(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		uuid uuid.UUID
	}
	testUuid, _ := uuid.Parse("2b1b425e-dee2-4227-8d94-f470a0ce0cd0")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			fields: fields{
				db: nil,
			},
			args: args{
				uuid: testUuid,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM persons WHERE alt_id = \\$1").
					WithArgs("2b1b425e-dee2-4227-8d94-f470a0ce0cd0").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "DB error",
			fields: fields{
				db: nil,
			},
			args: args{
				uuid: testUuid,
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM persons WHERE alt_id = \\$1").
					WithArgs("2b1b425e-dee2-4227-8d94-f470a0ce0cd0").WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			results, _ := p.DeleteByUuid(tt.args.uuid)
			if tt.wantErr && results.Deleted != false {
				t.Errorf("PersonRepositoryImpl.DeleteByUuid() = %v, want %v", results.Deleted, true)
			}
		})
	}
}

func TestPersonRepositoryImpl_Update(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		request UpdatePersonRequest
	}
	testUuid, _ := uuid.Parse("2b1b425e-dee2-4227-8d94-f470a0ce0cd0")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		prepare func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Success",
			fields: fields{
				db: nil,
			},
			args: args{
				request: UpdatePersonRequest{
					Id:           testUuid,
					Name:         "test name",
					Email:        "test email",
					LastChangeBy: "test user",
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name")
				mock.ExpectQuery("UPDATE persons SET name = \\$1, email = \\$2, last_change_by = \\$3 WHERE alt_id = \\$4").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "DB error",
			fields: fields{
				db: nil,
			},
			args: args{
				request: UpdatePersonRequest{
					Id:           testUuid,
					Name:         "test name",
					Email:        "test email",
					LastChangeBy: "test user",
				},
			},
			prepare: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("UPDATE persons SET name = \\$1, email = \\$2, last_change_by = \\$3 WHERE alt_id = \\$4").
					WillReturnError(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)
			tt.prepare(mock)
			tt.fields.db = sqlx.NewDb(db, "sqlmock")
			p := &PersonRepositoryImpl{
				db: tt.fields.db,
			}
			_, err := p.Update(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepositoryImpl.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewPersonRepository(t *testing.T) {
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewPersonRepository(tt.args.db)
		})
	}
}
