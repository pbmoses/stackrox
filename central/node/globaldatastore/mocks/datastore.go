// Code generated by MockGen. DO NOT EDIT.
// Source: datastore.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datastore "github.com/stackrox/rox/central/node/datastore"
)

// MockGlobalDataStore is a mock of GlobalDataStore interface.
type MockGlobalDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalDataStoreMockRecorder
}

// MockGlobalDataStoreMockRecorder is the mock recorder for MockGlobalDataStore.
type MockGlobalDataStoreMockRecorder struct {
	mock *MockGlobalDataStore
}

// NewMockGlobalDataStore creates a new mock instance.
func NewMockGlobalDataStore(ctrl *gomock.Controller) *MockGlobalDataStore {
	mock := &MockGlobalDataStore{ctrl: ctrl}
	mock.recorder = &MockGlobalDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalDataStore) EXPECT() *MockGlobalDataStoreMockRecorder {
	return m.recorder
}

// CountAllNodes mocks base method.
func (m *MockGlobalDataStore) CountAllNodes(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAllNodes", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAllNodes indicates an expected call of CountAllNodes.
func (mr *MockGlobalDataStoreMockRecorder) CountAllNodes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAllNodes", reflect.TypeOf((*MockGlobalDataStore)(nil).CountAllNodes), ctx)
}

// GetAllClusterNodeStores mocks base method.
func (m *MockGlobalDataStore) GetAllClusterNodeStores(ctx context.Context, writeAccess bool) (map[string]datastore.DataStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllClusterNodeStores", ctx, writeAccess)
	ret0, _ := ret[0].(map[string]datastore.DataStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllClusterNodeStores indicates an expected call of GetAllClusterNodeStores.
func (mr *MockGlobalDataStoreMockRecorder) GetAllClusterNodeStores(ctx, writeAccess interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllClusterNodeStores", reflect.TypeOf((*MockGlobalDataStore)(nil).GetAllClusterNodeStores), ctx, writeAccess)
}

// GetClusterNodeStore mocks base method.
func (m *MockGlobalDataStore) GetClusterNodeStore(ctx context.Context, clusterID string, writeAccess bool) (datastore.DataStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterNodeStore", ctx, clusterID, writeAccess)
	ret0, _ := ret[0].(datastore.DataStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterNodeStore indicates an expected call of GetClusterNodeStore.
func (mr *MockGlobalDataStoreMockRecorder) GetClusterNodeStore(ctx, clusterID, writeAccess interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterNodeStore", reflect.TypeOf((*MockGlobalDataStore)(nil).GetClusterNodeStore), ctx, clusterID, writeAccess)
}

// RemoveClusterNodeStores mocks base method.
func (m *MockGlobalDataStore) RemoveClusterNodeStores(ctx context.Context, clusterIDs ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range clusterIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveClusterNodeStores", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveClusterNodeStores indicates an expected call of RemoveClusterNodeStores.
func (mr *MockGlobalDataStoreMockRecorder) RemoveClusterNodeStores(ctx interface{}, clusterIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, clusterIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveClusterNodeStores", reflect.TypeOf((*MockGlobalDataStore)(nil).RemoveClusterNodeStores), varargs...)
}
