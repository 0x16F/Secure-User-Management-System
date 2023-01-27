// Code generated by MockGen. DO NOT EDIT.
// Source: model.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	reflect "reflect"
	user "test-project/src/internal/user"

	gomock "github.com/golang/mock/gomock"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStorager) Create(dto *user.UserDTO) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", dto)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockStoragerMockRecorder) Create(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStorager)(nil).Create), dto)
}

// Delete mocks base method.
func (m *MockStorager) Delete(id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStoragerMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorager)(nil).Delete), id)
}

// FindAll mocks base method.
func (m *MockStorager) FindAll(limit, offset int) (*[]user.FindUserDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", limit, offset)
	ret0, _ := ret[0].(*[]user.FindUserDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockStoragerMockRecorder) FindAll(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockStorager)(nil).FindAll), limit, offset)
}

// FindByLogin mocks base method.
func (m *MockStorager) FindByLogin(login string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLogin", login)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLogin indicates an expected call of FindByLogin.
func (mr *MockStoragerMockRecorder) FindByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLogin", reflect.TypeOf((*MockStorager)(nil).FindByLogin), login)
}

// FindOne mocks base method.
func (m *MockStorager) FindOne(id int64) (*user.FindUserDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", id)
	ret0, _ := ret[0].(*user.FindUserDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockStoragerMockRecorder) FindOne(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockStorager)(nil).FindOne), id)
}

// Update mocks base method.
func (m *MockStorager) Update(id int64, dto *user.UpdateUserDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockStoragerMockRecorder) Update(id, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorager)(nil).Update), id, dto)
}
