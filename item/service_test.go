package item

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateItem(t *testing.T) {
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

func TestUpdateItem(t *testing.T) {
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
