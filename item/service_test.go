package item

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"testing"
)

func TestItemService_CreateItem(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockItemRepository(controller)
	service := NewItemService(mockRepo)

	tests := []struct {
		name            string
		givenRequest    CreateItemRequest
		mockReturnValue ItemRow
		mockError       error
		expectedError   error
	}{
		{
			name:            "ValidRequest",
			givenRequest:    CreateItemRequest{Name: "item1", Description: "description1", UnitPrice: 100.00, CreatedBy: "testuser"},
			mockReturnValue: ItemRow{Id: 1, AltId: uuid.New(), Name: "item1", Description: "description1", UnitPrice: 100.00, CreatedBy: "testuser"},
			mockError:       nil,
			expectedError:   nil,
		},
		{
			name:            "RepoError",
			givenRequest:    CreateItemRequest{Name: "item1", Description: "description1", UnitPrice: 100.00, CreatedBy: "testuser"},
			mockReturnValue: ItemRow{},
			mockError:       errors.New("DB Error"),
			expectedError:   errors.New("DB Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().CreateItem(tt.givenRequest).Return(tt.mockReturnValue, tt.mockError)

			newItem, err := service.CreateItem(tt.givenRequest)

			if tt.mockError != nil {
				assert.Nil(t, newItem)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				mockValue := itemFromRow(tt.mockReturnValue)
				assert.Equal(t, &mockValue, newItem)
				assert.Nil(t, err)
			}
		})
	}
}

func TestItemService_UpdateItem(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockItemRepository(controller)
	service := NewItemService(mockRepo)

	tests := []struct {
		name            string
		givenRequest    UpdateItemRequest
		mockReturnValue ItemRow
		mockError       error
		expectedError   error
	}{
		{
			name:            "ValidRequest",
			givenRequest:    UpdateItemRequest{Id: uuid.New(), Name: "item1", Description: "description1", UnitPrice: 100.00, LastChangedBy: "testuser"},
			mockReturnValue: ItemRow{Id: 1, AltId: uuid.New(), Name: "item1", Description: "description1", UnitPrice: 100.00, CreatedBy: "testuser"},
			mockError:       nil,
			expectedError:   nil,
		},
		{
			name:            "RepoError",
			givenRequest:    UpdateItemRequest{Id: uuid.New(), Name: "item1", Description: "description1", UnitPrice: 100.00, LastChangedBy: "testuser"},
			mockReturnValue: ItemRow{},
			mockError:       errors.New("DB Error"),
			expectedError:   errors.New("DB Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().UpdateItem(tt.givenRequest).Return(tt.mockReturnValue, tt.mockError)

			newItem, err := service.UpdateItem(tt.givenRequest)

			if tt.mockError != nil {
				assert.Nil(t, newItem)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				mockValue := itemFromRow(tt.mockReturnValue)
				assert.Equal(t, &mockValue, newItem)
				assert.Nil(t, err)
			}
		})
	}
}

func TestItemService_DeleteItem(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockItemRepository(controller)
	service := NewItemService(mockRepo)

	tests := []struct {
		name            string
		givenId         uuid.UUID
		mockReturnValue commons.DeleteResult
		mockError       error
		expectedError   error
	}{
		{
			name:            "ValidRequest",
			givenId:         uuid.New(),
			mockReturnValue: commons.DeleteResult{Deleted: true},
			mockError:       nil,
			expectedError:   nil,
		},
		{
			name:            "RepoError",
			givenId:         uuid.New(),
			mockReturnValue: commons.DeleteResult{},
			mockError:       errors.New("DB Error"),
			expectedError:   errors.New("DB Error"),
		},
		{
			name:            "NonExistingId",
			givenId:         uuid.New(),
			mockReturnValue: commons.DeleteResult{Deleted: false},
			mockError:       nil,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().DeleteItem(tt.givenId).Return(tt.mockReturnValue, tt.mockError)

			deleteResult, err := service.DeleteItem(tt.givenId)

			if tt.mockError != nil {
				assert.Nil(t, deleteResult)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.Equal(t, tt.mockReturnValue.Deleted, deleteResult.Deleted)
				assert.Nil(t, err)
			}
		})
	}
}

func TestItemService_GetItem(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockRepo := NewMockItemRepository(controller)
	service := NewItemService(mockRepo)

	tests := []struct {
		name            string
		givenId         uuid.UUID
		mockReturnValue ItemRow
		mockError       error
		expectedError   error
	}{
		{
			name:            "ValidRequest",
			givenId:         uuid.New(),
			mockReturnValue: ItemRow{Id: 1, AltId: uuid.New(), Name: "item1", Description: "description1", UnitPrice: 100.00, CreatedBy: "testuser"},
			mockError:       nil,
			expectedError:   nil,
		},
		{
			name:            "RepoError",
			givenId:         uuid.New(),
			mockReturnValue: ItemRow{},
			mockError:       errors.New("DB Error"),
			expectedError:   errors.New("DB Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetItem(tt.givenId).Return(tt.mockReturnValue, tt.mockError)

			newItem, err := service.GetItem(tt.givenId)

			if tt.mockError != nil {
				assert.Nil(t, newItem)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				mockValue := itemFromRow(tt.mockReturnValue)
				assert.Equal(t, &mockValue, newItem)
				assert.Nil(t, err)
			}
		})
	}
}
