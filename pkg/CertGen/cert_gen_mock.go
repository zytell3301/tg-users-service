// Code generated by MockGen. DO NOT EDIT.
// Source: cert_gen_interface.go

// Package CertGen is a generated GoMock package.
package CertGen

import (
	x509 "crypto/x509"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGen is a mock of Gen interface.
type MockGen struct {
	ctrl     *gomock.Controller
	recorder *MockGenMockRecorder
}

// MockGenMockRecorder is the mock recorder for MockGen.
type MockGenMockRecorder struct {
	mock *MockGen
}

// NewMockGen creates a new mock instance.
func NewMockGen(ctrl *gomock.Controller) *MockGen {
	mock := &MockGen{ctrl: ctrl}
	mock.recorder = &MockGenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGen) EXPECT() *MockGenMockRecorder {
	return m.recorder
}

// NewCertificate mocks base method.
func (m *MockGen) NewCertificate(cert *x509.Certificate) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCertificate", cert)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCertificate indicates an expected call of NewCertificate.
func (mr *MockGenMockRecorder) NewCertificate(cert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCertificate", reflect.TypeOf((*MockGen)(nil).NewCertificate), cert)
}
