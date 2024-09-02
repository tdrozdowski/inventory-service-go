package person

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"inventory-service-go/commons"
	"time"
)

// PersonRow struct
type PersonRow struct {
	Id           int       `db:"id"`
	AltId        uuid.UUID `db:"alt_id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	CreatedBy    string    `db:"created_by"`
	CreatedAt    time.Time `db:"created_at"`
	LastUpdate   time.Time `db:"last_update"`
	LastChangeBy string    `db:"last_change_by"`
}

type CreatePersonRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedBy string `json:"created_by"`
}

type UpdatePersonRequest struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	LastChangeBy string    `json:"last_change_by"`
}

// Interface for PersonRepository
type PersonRepository interface {
	GetAll(pagination *commons.Pagination) ([]PersonRow, error)
	GetById(id int) (PersonRow, error)
	GetByUuid(uuid uuid.UUID) (PersonRow, error)
	Create(request CreatePersonRequest) (PersonRow, error)
	Update(request UpdatePersonRequest) (PersonRow, error)
	DeleteByUuid(uuid uuid.UUID) (commons.DeleteResult, error)
}

type PersonRepositoryImpl struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) PersonRepository {
	return &PersonRepositoryImpl{
		db: db,
	}
}

func (p *PersonRepositoryImpl) GetAll(pagination *commons.Pagination) ([]PersonRow, error) {
	// uses sqlx to query the perons table and retreive all rows
	var persons []PersonRow
	if pagination != nil {
		err := p.db.Select(&persons, "SELECT * FROM persons WHERE id > ? LIMIT ?", pagination.LastId, pagination.PageSize)
		if err != nil {
			return nil, err
		}
	} else {
		err := p.db.Select(&persons, "SELECT * FROM persons")
		if err != nil {
			return nil, err
		}
	}
	return persons, nil
}

func (p *PersonRepositoryImpl) GetById(id int) (PersonRow, error) {
	// uses sqlx to query the persons table and retrieve a single row by id
	var person PersonRow
	err := p.db.Get(&person, "SELECT * FROM persons WHERE id = ?", id)
	if err != nil {
		return PersonRow{}, err
	}
	return person, nil
}

func (p *PersonRepositoryImpl) GetByUuid(uuid uuid.UUID) (PersonRow, error) {
	// uses sqlx to query the persons table and retrieve a single row by uuid
	var person PersonRow
	err := p.db.Get(&person, "SELECT * FROM persons WHERE alt_id = ?", uuid)
	if err != nil {
		return PersonRow{}, err
	}
	return person, nil
}

func (p *PersonRepositoryImpl) Create(request CreatePersonRequest) (PersonRow, error) {
	// uses sqlx to insert a new row into the persons table
	var person PersonRow
	err := p.db.Get(&person, "INSERT INTO persons (name, email, created_by) VALUES (?, ?, ?) RETURNING *", request.Name, request.Email, request.CreatedBy)
	if err != nil {
		return PersonRow{}, err
	}
	return person, nil
}

func (p *PersonRepositoryImpl) Update(request UpdatePersonRequest) (PersonRow, error) {
	// uses sqlx to update a row in the persons table
	var person PersonRow
	err := p.db.Get(&person, "UPDATE persons SET name = ?, email = ?, last_change_by = ? WHERE alt_id = ? RETURNING *", request.Name, request.Email, request.LastChangeBy, request.Id)
	if err != nil {
		return PersonRow{}, err
	}
	return person, nil
}

func (p *PersonRepositoryImpl) DeleteByUuid(uuid uuid.UUID) (commons.DeleteResult, error) {
	// uses sqlx to delete a row from the persons table by uuid
	sqlResults := p.db.MustExec("DELETE FROM persons WHERE alt_id = ?", uuid)
	rowsAffected, err := sqlResults.RowsAffected()
	if err != nil {
		return commons.DeleteResult{}, err
	}
	result := commons.DeleteResult{
		Id:      uuid,
		Deleted: rowsAffected > 0,
	}
	return result, nil
}
