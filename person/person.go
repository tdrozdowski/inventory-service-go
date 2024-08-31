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

// Interface for PersonRepository
type PersonRepository interface {
	GetAll() ([]PersonRow, error)
	GetById(id int) (PersonRow, error)
	Create(person PersonRow) (PersonRow, error)
	Update(person PersonRow) (PersonRow, error)
	Delete(id int) (commons.DeleteResult, error)
}

type PersonRepositoryImpl struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) PersonRepository {
	return &PersonRepositoryImpl{
		db: db,
	}
}

func (p *PersonRepositoryImpl) GetAll() ([]PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) GetById(id int) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Create(person PersonRow) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Update(person PersonRow) (PersonRow, error) {
	panic("implement me")
}

func (p *PersonRepositoryImpl) Delete(id int) (commons.DeleteResult, error) {
	panic("implement me")
}
