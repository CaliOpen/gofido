// Code generated by MockGen. DO NOT EDIT.
// Source: store/cassandra.go

// Package store is a generated GoMock package.
package store

import (
	gomock "github.com/golang/mock/gomock"
	u2f "github.com/tstranex/u2f"
	reflect "reflect"
)

// MockStoreInterface is a mock of StoreInterface interface
type MockStoreInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStoreInterfaceMockRecorder
}

// MockStoreInterfaceMockRecorder is the mock recorder for MockStoreInterface
type MockStoreInterfaceMockRecorder struct {
	mock *MockStoreInterface
}

// NewMockStoreInterface creates a new mock instance
func NewMockStoreInterface(ctrl *gomock.Controller) *MockStoreInterface {
	mock := &MockStoreInterface{ctrl: ctrl}
	mock.recorder = &MockStoreInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStoreInterface) EXPECT() *MockStoreInterfaceMockRecorder {
	return m.recorder
}

// NewChallenge mocks base method
func (m *MockStoreInterface) NewChallenge(arg0 string) (u2f.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewChallenge", arg0)
	ret0, _ := ret[0].(u2f.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewChallenge indicates an expected call of NewChallenge
func (mr *MockStoreInterfaceMockRecorder) NewChallenge(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewChallenge", reflect.TypeOf((*MockStoreInterface)(nil).NewChallenge), arg0)
}

// GetChallenge mocks base method
func (m *MockStoreInterface) GetChallenge(arg0, arg1 string) (u2f.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallenge", arg0, arg1)
	ret0, _ := ret[0].(u2f.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallenge indicates an expected call of GetChallenge
func (mr *MockStoreInterfaceMockRecorder) GetChallenge(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallenge", reflect.TypeOf((*MockStoreInterface)(nil).GetChallenge), arg0, arg1)
}

// NewRegistration mocks base method
func (m *MockStoreInterface) NewRegistration(arg0 string, arg1 u2f.Challenge, arg2 u2f.RegisterResponse) (*u2f.Registration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRegistration", arg0, arg1, arg2)
	ret0, _ := ret[0].(*u2f.Registration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRegistration indicates an expected call of NewRegistration
func (mr *MockStoreInterfaceMockRecorder) NewRegistration(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRegistration", reflect.TypeOf((*MockStoreInterface)(nil).NewRegistration), arg0, arg1, arg2)
}

// GetRegistrations mocks base method
func (m *MockStoreInterface) GetRegistrations(arg0 string) ([]u2f.Registration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegistrations", arg0)
	ret0, _ := ret[0].([]u2f.Registration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegistrations indicates an expected call of GetRegistrations
func (mr *MockStoreInterfaceMockRecorder) GetRegistrations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegistrations", reflect.TypeOf((*MockStoreInterface)(nil).GetRegistrations), arg0)
}

// InsertKeyChallenge mocks base method
func (m *MockStoreInterface) InsertKeyChallenge(arg0 string, arg1 []byte, arg2 u2f.Challenge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertKeyChallenge", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertKeyChallenge indicates an expected call of InsertKeyChallenge
func (mr *MockStoreInterfaceMockRecorder) InsertKeyChallenge(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertKeyChallenge", reflect.TypeOf((*MockStoreInterface)(nil).InsertKeyChallenge), arg0, arg1, arg2)
}

// GetKeyChallenges mocks base method
func (m *MockStoreInterface) GetKeyChallenges(arg0 string, arg1 []byte) []KeyChallenge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeyChallenges", arg0, arg1)
	ret0, _ := ret[0].([]KeyChallenge)
	return ret0
}

// GetKeyChallenges indicates an expected call of GetKeyChallenges
func (mr *MockStoreInterfaceMockRecorder) GetKeyChallenges(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeyChallenges", reflect.TypeOf((*MockStoreInterface)(nil).GetKeyChallenges), arg0, arg1)
}

// GetKeyCounter mocks base method
func (m *MockStoreInterface) GetKeyCounter(arg0 string, arg1 []byte) (KeyCounter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeyCounter", arg0, arg1)
	ret0, _ := ret[0].(KeyCounter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeyCounter indicates an expected call of GetKeyCounter
func (mr *MockStoreInterfaceMockRecorder) GetKeyCounter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeyCounter", reflect.TypeOf((*MockStoreInterface)(nil).GetKeyCounter), arg0, arg1)
}

// UpdateCounter mocks base method
func (m *MockStoreInterface) UpdateCounter(arg0 string, arg1 []byte, arg2 uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter
func (mr *MockStoreInterfaceMockRecorder) UpdateCounter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockStoreInterface)(nil).UpdateCounter), arg0, arg1, arg2)
}