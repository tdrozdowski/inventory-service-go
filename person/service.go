package person

import (
	"github.com/google/uuid"
	"inventory-service-go/commons"
)

type Person struct {
	Seq       int               `json:"seq"`
	Id        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	AuditInfo commons.AuditInfo `json:"audit_info"`
}

func (*Person) FromRow(row PersonRow) Person {
	return Person{
		Seq:   row.Id,
		Id:    row.AltId,
		Name:  row.Name,
		Email: row.Email,
		AuditInfo: commons.AuditInfo{
			CreatedBy:     row.CreatedBy,
			CreatedAt:     row.CreatedAt.String(),
			LastUpdate:    row.LastUpdate.String(),
			LastChangedBy: row.LastChangedBy,
		},
	}
}

type PersonService interface {
	GetAll(pagination *commons.Pagination) ([]Person, error)
	GetById(id uuid.UUID) (Person, error)
	Create(request CreatePersonRequest) (Person, error)
	Update(request UpdatePersonRequest) (Person, error)
	DeleteByUuid(uuid uuid.UUID) (commons.DeleteResult, error)
}

type PersonServiceImpl struct {
	repo PersonRepository
}

func NewPersonService(repo PersonRepository) PersonServiceImpl {
	return PersonServiceImpl{
		repo: repo,
	}
}

func (p PersonServiceImpl) GetAll(pagination *commons.Pagination) ([]Person, error) {
	persons, err := p.repo.GetAll(pagination)
	if err != nil {
		return nil, err
	}

	var result []Person
	for _, person := range persons {
		p2 := Person{}
		result = append(result, p2.FromRow(person))
	}

	return result, nil
}

func (p PersonServiceImpl) GetById(id uuid.UUID) (Person, error) {
	person, err := p.repo.GetByUuid(id)
	p2 := Person{}
	if err != nil {
		return p2, err
	}

	return p2.FromRow(person), nil
}

func (p PersonServiceImpl) Create(request CreatePersonRequest) (Person, error) {
	person, err := p.repo.Create(request)
	p2 := Person{}
	if err != nil {
		return p2, err
	}

	return p2.FromRow(person), nil
}

func (p PersonServiceImpl) Update(request UpdatePersonRequest) (Person, error) {
	person, err := p.repo.Update(request)
	p2 := Person{}
	if err != nil {
		return p2, err
	}

	return p2.FromRow(person), nil
}

func (p PersonServiceImpl) DeleteByUuid(uuid uuid.UUID) (commons.DeleteResult, error) {
	return p.repo.DeleteByUuid(uuid)
}
