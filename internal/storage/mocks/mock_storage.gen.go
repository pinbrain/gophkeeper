// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pinbrain/gophkeeper/internal/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// CreateItem mocks base method.
func (m *MockStorage) CreateItem(ctx context.Context, userID string, item *model.VaultItem) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateItem", ctx, userID, item)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateItem indicates an expected call of CreateItem.
func (mr *MockStorageMockRecorder) CreateItem(ctx, userID, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateItem", reflect.TypeOf((*MockStorage)(nil).CreateItem), ctx, userID, item)
}

// CreateUser mocks base method.
func (m *MockStorage) CreateUser(ctx context.Context, user *model.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStorageMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorage)(nil).CreateUser), ctx, user)
}

// DeleteItem mocks base method.
func (m *MockStorage) DeleteItem(ctx context.Context, id, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockStorageMockRecorder) DeleteItem(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockStorage)(nil).DeleteItem), ctx, id, userID)
}

// GetItem mocks base method.
func (m *MockStorage) GetItem(ctx context.Context, id, userID string) (*model.VaultItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, id, userID)
	ret0, _ := ret[0].(*model.VaultItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockStorageMockRecorder) GetItem(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockStorage)(nil).GetItem), ctx, id, userID)
}

// GetItemsByType mocks base method.
func (m *MockStorage) GetItemsByType(ctx context.Context, dataType, userID string) ([]model.VaultItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByType", ctx, dataType, userID)
	ret0, _ := ret[0].([]model.VaultItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsByType indicates an expected call of GetItemsByType.
func (mr *MockStorageMockRecorder) GetItemsByType(ctx, dataType, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByType", reflect.TypeOf((*MockStorage)(nil).GetItemsByType), ctx, dataType, userID)
}

// GetUserByID mocks base method.
func (m *MockStorage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStorageMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStorage)(nil).GetUserByID), ctx, id)
}

// GetUserByLogin mocks base method.
func (m *MockStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", ctx, login)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStorageMockRecorder) GetUserByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorage)(nil).GetUserByLogin), ctx, login)
}

// UpdateItem mocks base method.
func (m *MockStorage) UpdateItem(ctx context.Context, id, userID string, item *model.VaultItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", ctx, id, userID, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockStorageMockRecorder) UpdateItem(ctx, id, userID, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockStorage)(nil).UpdateItem), ctx, id, userID, item)
}

// MockUserStorage is a mock of UserStorage interface.
type MockUserStorage struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorageMockRecorder
}

// MockUserStorageMockRecorder is the mock recorder for MockUserStorage.
type MockUserStorageMockRecorder struct {
	mock *MockUserStorage
}

// NewMockUserStorage creates a new mock instance.
func NewMockUserStorage(ctrl *gomock.Controller) *MockUserStorage {
	mock := &MockUserStorage{ctrl: ctrl}
	mock.recorder = &MockUserStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorage) EXPECT() *MockUserStorageMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserStorage) CreateUser(ctx context.Context, user *model.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserStorageMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserStorage)(nil).CreateUser), ctx, user)
}

// GetUserByID mocks base method.
func (m *MockUserStorage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserStorageMockRecorder) GetUserByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserStorage)(nil).GetUserByID), ctx, id)
}

// GetUserByLogin mocks base method.
func (m *MockUserStorage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", ctx, login)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserStorageMockRecorder) GetUserByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserStorage)(nil).GetUserByLogin), ctx, login)
}

// MockVaultStorage is a mock of VaultStorage interface.
type MockVaultStorage struct {
	ctrl     *gomock.Controller
	recorder *MockVaultStorageMockRecorder
}

// MockVaultStorageMockRecorder is the mock recorder for MockVaultStorage.
type MockVaultStorageMockRecorder struct {
	mock *MockVaultStorage
}

// NewMockVaultStorage creates a new mock instance.
func NewMockVaultStorage(ctrl *gomock.Controller) *MockVaultStorage {
	mock := &MockVaultStorage{ctrl: ctrl}
	mock.recorder = &MockVaultStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVaultStorage) EXPECT() *MockVaultStorageMockRecorder {
	return m.recorder
}

// CreateItem mocks base method.
func (m *MockVaultStorage) CreateItem(ctx context.Context, userID string, item *model.VaultItem) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateItem", ctx, userID, item)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateItem indicates an expected call of CreateItem.
func (mr *MockVaultStorageMockRecorder) CreateItem(ctx, userID, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateItem", reflect.TypeOf((*MockVaultStorage)(nil).CreateItem), ctx, userID, item)
}

// DeleteItem mocks base method.
func (m *MockVaultStorage) DeleteItem(ctx context.Context, id, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockVaultStorageMockRecorder) DeleteItem(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockVaultStorage)(nil).DeleteItem), ctx, id, userID)
}

// GetItem mocks base method.
func (m *MockVaultStorage) GetItem(ctx context.Context, id, userID string) (*model.VaultItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, id, userID)
	ret0, _ := ret[0].(*model.VaultItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockVaultStorageMockRecorder) GetItem(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockVaultStorage)(nil).GetItem), ctx, id, userID)
}

// GetItemsByType mocks base method.
func (m *MockVaultStorage) GetItemsByType(ctx context.Context, dataType, userID string) ([]model.VaultItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByType", ctx, dataType, userID)
	ret0, _ := ret[0].([]model.VaultItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsByType indicates an expected call of GetItemsByType.
func (mr *MockVaultStorageMockRecorder) GetItemsByType(ctx, dataType, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByType", reflect.TypeOf((*MockVaultStorage)(nil).GetItemsByType), ctx, dataType, userID)
}

// UpdateItem mocks base method.
func (m *MockVaultStorage) UpdateItem(ctx context.Context, id, userID string, item *model.VaultItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", ctx, id, userID, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockVaultStorageMockRecorder) UpdateItem(ctx, id, userID, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockVaultStorage)(nil).UpdateItem), ctx, id, userID, item)
}
