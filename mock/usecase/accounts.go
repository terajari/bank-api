// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/accounts.go

// Package mockusecase is a generated GoMock package.
package mockusecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/terajari/bank-api/dto"
)

// MockAccountsUsecase is a mock of AccountsUsecase interface.
type MockAccountsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAccountsUsecaseMockRecorder
}

// MockAccountsUsecaseMockRecorder is the mock recorder for MockAccountsUsecase.
type MockAccountsUsecaseMockRecorder struct {
	mock *MockAccountsUsecase
}

// NewMockAccountsUsecase creates a new mock instance.
func NewMockAccountsUsecase(ctrl *gomock.Controller) *MockAccountsUsecase {
	mock := &MockAccountsUsecase{ctrl: ctrl}
	mock.recorder = &MockAccountsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountsUsecase) EXPECT() *MockAccountsUsecaseMockRecorder {
	return m.recorder
}

// DeleteAccount mocks base method.
func (m *MockAccountsUsecase) DeleteAccount(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockAccountsUsecaseMockRecorder) DeleteAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockAccountsUsecase)(nil).DeleteAccount), ctx, id)
}

// GetAccount mocks base method.
func (m *MockAccountsUsecase) GetAccount(ctx context.Context, id string) (dto.GetAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, id)
	ret0, _ := ret[0].(dto.GetAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAccountsUsecaseMockRecorder) GetAccount(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountsUsecase)(nil).GetAccount), ctx, id)
}

// ListAccounts mocks base method.
func (m *MockAccountsUsecase) ListAccounts(ctx context.Context, req dto.ListAccountsRequest) ([]dto.GetAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", ctx, req)
	ret0, _ := ret[0].([]dto.GetAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockAccountsUsecaseMockRecorder) ListAccounts(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockAccountsUsecase)(nil).ListAccounts), ctx, req)
}

// RegisterNewAccounts mocks base method.
func (m *MockAccountsUsecase) RegisterNewAccounts(ctx context.Context, req dto.RegisterNewAccountsRequest) (dto.RegisterNewAccountsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNewAccounts", ctx, req)
	ret0, _ := ret[0].(dto.RegisterNewAccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterNewAccounts indicates an expected call of RegisterNewAccounts.
func (mr *MockAccountsUsecaseMockRecorder) RegisterNewAccounts(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNewAccounts", reflect.TypeOf((*MockAccountsUsecase)(nil).RegisterNewAccounts), ctx, req)
}

// UpdateAccount mocks base method.
func (m *MockAccountsUsecase) UpdateAccount(ctx context.Context, req dto.UpdateAccountRequest) (dto.UpdateAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, req)
	ret0, _ := ret[0].(dto.UpdateAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockAccountsUsecaseMockRecorder) UpdateAccount(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockAccountsUsecase)(nil).UpdateAccount), ctx, req)
}