// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zytell3301/tg-error-reporter (interfaces: Reporter)

// Package mock_tg_error_reporter is a generated GoMock package.
package core

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ErrorReporter "github.com/zytell3301/tg-error-reporter"
)

// MockReporter is a mock of Reporter interface.
type MockReporter struct {
	ctrl     *gomock.Controller
	recorder *MockReporterMockRecorder
}

// MockReporterMockRecorder is the mock recorder for MockReporter.
type MockReporterMockRecorder struct {
	mock *MockReporter
}

// NewMockReporter creates a new mock instance.
func NewMockReporter(ctrl *gomock.Controller) *MockReporter {
	mock := &MockReporter{ctrl: ctrl}
	mock.recorder = &MockReporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReporter) EXPECT() *MockReporterMockRecorder {
	return m.recorder
}

// Report mocks base method.
func (m *MockReporter) Report(arg0 ErrorReporter.Error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Report", arg0)
}

// Report indicates an expected call of Report.
func (mr *MockReporterMockRecorder) Report(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Report", reflect.TypeOf((*MockReporter)(nil).Report), arg0)
}