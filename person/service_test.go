package person

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"reflect"
	"testing"
	"time"
)

func TestFromRow(t *testing.T) {
	altId := uuid.New()
	now := time.Now()
	tests := []struct {
		name  string
		input PersonRow
		want  Person
	}{
		{
			name: "Complete Person Information",
			input: PersonRow{
				Id:            1,
				AltId:         altId,
				Name:          "John Doe",
				Email:         "johndoe@example.com",
				CreatedBy:     "admin",
				CreatedAt:     now,
				LastUpdate:    now,
				LastChangedBy: "admin",
			},
			want: Person{
				Seq:   1,
				Id:    altId,
				Name:  "John Doe",
				Email: "johndoe@example.com",
				AuditInfo: commons.AuditInfo{
					CreatedBy:     "admin",
					CreatedAt:     now.String(),
					LastUpdate:    now.String(),
					LastChangedBy: "admin",
				},
			},
		},
		// add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Person{}
			got := s.FromRow(tt.input)
			if got != tt.want {
				t.Errorf("Person.FromRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func personRowFixture() PersonRow {
	return PersonRow{
		Id:            1,
		AltId:         uuid.New(),
		Name:          "John Doe",
		Email:         "john.doe@test.com",
		CreatedBy:     "unit_test",
		CreatedAt:     time.Now(),
		LastUpdate:    time.Now(),
		LastChangedBy: "unit_test",
	}
}

func personFixture(personRow PersonRow) Person {
	p := &Person{}
	return p.FromRow(personRow)
}

func TestGetAll(t *testing.T) {
	pagination := &commons.Pagination{LastId: 0, PageSize: 10}
	controller := gomock.NewController(t)
	mockRepo := NewMockPersonRepository(controller)
	personService := NewPersonService(mockRepo)
	rowFixture := personRowFixture()
	expectedPerson := personFixture(rowFixture)
	tests := []struct {
		name     string
		expected []PersonRow
		want     []Person
		wantErr  bool
	}{
		{
			name: "Get All Persons",
			expected: []PersonRow{
				rowFixture,
			},
			want: []Person{
				expectedPerson,
			},
			wantErr: false,
		},
		{
			name:     "Error on GetAll",
			expected: nil,
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		if tt.wantErr {
			mockRepo.EXPECT().GetAll(pagination).Return(nil, errors.New("error"))
		} else {
			mockRepo.EXPECT().GetAll(pagination).Return(tt.expected, nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := personService.GetAll(pagination)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersonServiceImpl.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PersonServiceImpl.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockPersonRepository(controller)
	personService := NewPersonService(mockRepo)
	rowFixture := personRowFixture()
	expectedPerson := personFixture(rowFixture)

	tests := []struct {
		name     string
		expected PersonRow
		want     Person
		wantErr  bool
	}{
		{
			name:     "Get Person By Id",
			expected: rowFixture,
			want:     expectedPerson,
			wantErr:  false,
		},
		{
			name:     "Error on GetById",
			expected: PersonRow{},
			want:     Person{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		if tt.wantErr {
			mockRepo.EXPECT().GetByUuid(tt.expected.AltId).Return(PersonRow{}, errors.New("error"))
		} else {
			mockRepo.EXPECT().GetByUuid(tt.expected.AltId).Return(tt.expected, nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := personService.GetById(tt.expected.AltId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersonServiceImpl.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("PersonServiceImpl.GetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockPersonRepository(controller)
	personService := NewPersonService(mockRepo)
	rowFixture := personRowFixture()
	expectedPerson := personFixture(rowFixture)
	createRequest := CreatePersonRequest{
		Name:      "John Doe",
		Email:     "john.doe@test.com",
		CreatedBy: "unit_test",
	}
	tests := []struct {
		name     string
		expected PersonRow
		want     Person
		wantErr  bool
	}{
		{
			name:     "Create Person",
			expected: rowFixture,
			want:     expectedPerson,
			wantErr:  false,
		},
		{
			name:     "Error on Create",
			expected: PersonRow{},
			want:     Person{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		if tt.wantErr {
			mockRepo.EXPECT().Create(createRequest).Return(PersonRow{}, errors.New("error"))
		} else {
			mockRepo.EXPECT().Create(createRequest).Return(tt.expected, nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := personService.Create(createRequest)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersonServiceImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("PersonServiceImpl.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockPersonRepository(controller)
	personService := NewPersonService(mockRepo)
	rowFixture := personRowFixture()
	expectedPerson := personFixture(rowFixture)
	updateRequest := UpdatePersonRequest{
		Id:           rowFixture.AltId,
		Name:         "John Doe",
		Email:        "john.doe@test.com",
		LastChangeBy: "unit_test",
	}
	tests := []struct {
		name     string
		expected PersonRow
		want     Person
		wantErr  bool
	}{
		{
			name:     "Update Person",
			expected: rowFixture,
			want:     expectedPerson,
			wantErr:  false,
		},
		{
			name:     "Error on Update",
			expected: PersonRow{},
			want:     Person{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		if tt.wantErr {
			mockRepo.EXPECT().Update(updateRequest).Return(PersonRow{}, errors.New("error"))
		} else {
			mockRepo.EXPECT().Update(updateRequest).Return(tt.expected, nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := personService.Update(updateRequest)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersonServiceImpl.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("PersonServiceImpl.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteByUuid(t *testing.T) {
	rowFixture := personRowFixture()
	tests := []struct {
		name     string
		expected commons.DeleteResult
		want     commons.DeleteResult
		wantErr  bool
	}{
		{
			name:     "Delete Person",
			expected: commons.DeleteResult{Id: rowFixture.AltId, Deleted: true},
			want:     commons.DeleteResult{Id: rowFixture.AltId, Deleted: true},
			wantErr:  false,
		},
		{
			name:     "Error on Delete",
			expected: commons.DeleteResult{},
			want:     commons.DeleteResult{},
			wantErr:  true,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockRepo := NewMockPersonRepository(controller)
		personService := NewPersonService(mockRepo)
		if tt.wantErr {
			mockRepo.EXPECT().DeleteByUuid(rowFixture.AltId).Return(commons.DeleteResult{}, errors.New("error"))
		} else {
			mockRepo.EXPECT().DeleteByUuid(rowFixture.AltId).Return(tt.expected, nil)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := personService.DeleteByUuid(rowFixture.AltId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PersonServiceImpl.DeleteByUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && err != nil && tt.wantErr {
				assert.True(t, true, "Error is expected")
			} else {
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("PersonServiceImpl.DeleteByUuid() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
