// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//      mockgen -source=client.go destination=osrmtest/mock_client.go
//
// Package mock_osrm is a generated GoMock package.
package mock_osrm

import (
        reflect "reflect"

        route "github.com/nextmv-io/sdk/route"
        osrm "github.com/nextmv-io/sdk/route/osrm"
        gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
        ctrl     *gomock.Controller
        recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
        mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
        mock := &MockClient{ctrl: ctrl}
        mock.recorder = &MockClientMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
        return m.recorder
}

// Get mocks base method.
func (m *MockClient) Get(uri string) ([]byte, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Get", uri)
        ret0, _ := ret[0].([]byte)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(uri any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), uri)
}

// IgnoreEmpty mocks base method.
func (m *MockClient) IgnoreEmpty(ignore bool) {
        m.ctrl.T.Helper()
        m.ctrl.Call(m, "IgnoreEmpty", ignore)
}

// IgnoreEmpty indicates an expected call of IgnoreEmpty.
func (mr *MockClientMockRecorder) IgnoreEmpty(ignore any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IgnoreEmpty", reflect.TypeOf((*MockClient)(nil).IgnoreEmpty), ignore)
}

// MaxTableSize mocks base method.
func (m *MockClient) MaxTableSize(size int) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "MaxTableSize", size)
        ret0, _ := ret[0].(error)
        return ret0
}

// MaxTableSize indicates an expected call of MaxTableSize.
func (mr *MockClientMockRecorder) MaxTableSize(size any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MaxTableSize", reflect.TypeOf((*MockClient)(nil).MaxTableSize), size)
}

// Polyline mocks base method.
func (m *MockClient) Polyline(points []route.Point) (string, []string, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Polyline", points)
        ret0, _ := ret[0].(string)
        ret1, _ := ret[1].([]string)
        ret2, _ := ret[2].(error)
        return ret0, ret1, ret2
}

// Polyline indicates an expected call of Polyline.
func (mr *MockClientMockRecorder) Polyline(points any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Polyline", reflect.TypeOf((*MockClient)(nil).Polyline), points)
}

// ScaleFactor mocks base method.
func (m *MockClient) ScaleFactor(factor float64) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "ScaleFactor", factor)
        ret0, _ := ret[0].(error)
        return ret0
}

// ScaleFactor indicates an expected call of ScaleFactor.
func (mr *MockClientMockRecorder) ScaleFactor(factor any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleFactor", reflect.TypeOf((*MockClient)(nil).ScaleFactor), factor)
}

// SnapRadius mocks base method.
func (m *MockClient) SnapRadius(radius int) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "SnapRadius", radius)
        ret0, _ := ret[0].(error)
        return ret0
}

// SnapRadius indicates an expected call of SnapRadius.
func (mr *MockClientMockRecorder) SnapRadius(radius any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SnapRadius", reflect.TypeOf((*MockClient)(nil).SnapRadius), radius)
}

// Table mocks base method.
func (m *MockClient) Table(points []route.Point, opts ...osrm.TableOptions) ([][]float64, [][]float64, error) {
        m.ctrl.T.Helper()
        varargs := []any{points}
        for _, a := range opts {
                varargs = append(varargs, a)
        }
        ret := m.ctrl.Call(m, "Table", varargs...)
        ret0, _ := ret[0].([][]float64)
        ret1, _ := ret[1].([][]float64)
        ret2, _ := ret[2].(error)
        return ret0, ret1, ret2
}

// Table indicates an expected call of Table.
func (mr *MockClientMockRecorder) Table(points any, opts ...any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        varargs := append([]any{points}, opts...)
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Table", reflect.TypeOf((*MockClient)(nil).Table), varargs...)
}