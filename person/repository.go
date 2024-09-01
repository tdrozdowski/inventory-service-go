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
	Delete(id int) (commons.DeleteResult, error)
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
	panic("implement me")
}

func (p *PersonRepositoryImpl) GetByUuid(uuid uuid.UUID) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Create(request CreatePersonRequest) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Update(request UpdatePersonRequest) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Delete(id int) (commons.DeleteResult, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) DeleteByUuid(uuid uuid.UUID) (commons.DeleteResult, error) {
	panic("implement me")
}
