// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source service.go -destination mock_service.go -package invoice
//

// Package invoice is a generated GoMock package.
package invoice

import (
	commons "inventory-service-go/commons"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockInvoiceService is a mock of InvoiceService interface.
type MockInvoiceService struct {
	ctrl     *gomock.Controller
	recorder *MockInvoiceServiceMockRecorder
}

// MockInvoiceServiceMockRecorder is the mock recorder for MockInvoiceService.
type MockInvoiceServiceMockRecorder struct {
	mock *MockInvoiceService
}

// NewMockInvoiceService creates a new mock instance.
func NewMockInvoiceService(ctrl *gomock.Controller) *MockInvoiceService {
	mock := &MockInvoiceService{ctrl: ctrl}
	mock.recorder = &MockInvoiceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInvoiceService) EXPECT() *MockInvoiceServiceMockRecorder {
	return m.recorder
}

// AddItemsToInvoice mocks base method.
func (m *MockInvoiceService) AddItemsToInvoice(request ItemsToInvoiceRequest) (ItemsToInvoiceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItemsToInvoice", request)
	ret0, _ := ret[0].(ItemsToInvoiceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddItemsToInvoice indicates an expected call of AddItemsToInvoice.
func (mr *MockInvoiceServiceMockRecorder) AddItemsToInvoice(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItemsToInvoice", reflect.TypeOf((*MockInvoiceService)(nil).AddItemsToInvoice), request)
}

// CreateInvoice mocks base method.
func (m *MockInvoiceService) CreateInvoice(invoice CreateInvoiceRequest) (Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInvoice", invoice)
	ret0, _ := ret[0].(Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateInvoice indicates an expected call of CreateInvoice.
func (mr *MockInvoiceServiceMockRecorder) CreateInvoice(invoice any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInvoice", reflect.TypeOf((*MockInvoiceService)(nil).CreateInvoice), invoice)
}

// DeleteInvoice mocks base method.
func (m *MockInvoiceService) DeleteInvoice(id uuid.UUID) (commons.DeleteResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteInvoice", id)
	ret0, _ := ret[0].(commons.DeleteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteInvoice indicates an expected call of DeleteInvoice.
func (mr *MockInvoiceServiceMockRecorder) DeleteInvoice(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteInvoice", reflect.TypeOf((*MockInvoiceService)(nil).DeleteInvoice), id)
}

// GetAllInvoices mocks base method.
func (m *MockInvoiceService) GetAllInvoices(pagination *commons.Pagination) ([]Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllInvoices", pagination)
	ret0, _ := ret[0].([]Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllInvoices indicates an expected call of GetAllInvoices.
func (mr *MockInvoiceServiceMockRecorder) GetAllInvoices(pagination any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllInvoices", reflect.TypeOf((*MockInvoiceService)(nil).GetAllInvoices), pagination)
}

// GetInvoice mocks base method.
func (m *MockInvoiceService) GetInvoice(id uuid.UUID, withItems bool) (Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoice", id, withItems)
	ret0, _ := ret[0].(Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoice indicates an expected call of GetInvoice.
func (mr *MockInvoiceServiceMockRecorder) GetInvoice(id, withItems any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoice", reflect.TypeOf((*MockInvoiceService)(nil).GetInvoice), id, withItems)
}

// GetInvoicesForUser mocks base method.
func (m *MockInvoiceService) GetInvoicesForUser(userId uuid.UUID) ([]Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoicesForUser", userId)
	ret0, _ := ret[0].([]Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvoicesForUser indicates an expected call of GetInvoicesForUser.
func (mr *MockInvoiceServiceMockRecorder) GetInvoicesForUser(userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoicesForUser", reflect.TypeOf((*MockInvoiceService)(nil).GetInvoicesForUser), userId)
}

// RemoveItemFromInvoice mocks base method.
func (m *MockInvoiceService) RemoveItemFromInvoice(request SimpleInvoiceItem) (ItemsToInvoiceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveItemFromInvoice", request)
	ret0, _ := ret[0].(ItemsToInvoiceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveItemFromInvoice indicates an expected call of RemoveItemFromInvoice.
func (mr *MockInvoiceServiceMockRecorder) RemoveItemFromInvoice(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveItemFromInvoice", reflect.TypeOf((*MockInvoiceService)(nil).RemoveItemFromInvoice), request)
}

// UpdateInvoice mocks base method.
func (m *MockInvoiceService) UpdateInvoice(invoice UpdateInvoiceRequest) (Invoice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInvoice", invoice)
	ret0, _ := ret[0].(Invoice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInvoice indicates an expected call of UpdateInvoice.
func (mr *MockInvoiceServiceMockRecorder) UpdateInvoice(invoice any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInvoice", reflect.TypeOf((*MockInvoiceService)(nil).UpdateInvoice), invoice)
}
