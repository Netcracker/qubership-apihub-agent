// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Netcracker/qubership-apihub-agent/client (interfaces: ApihubClient)

// Package mock_client is a generated GoMock package.
package mock_client

import (
	reflect "reflect"

	secctx "github.com/Netcracker/qubership-apihub-agent/secctx"
	view "github.com/Netcracker/qubership-apihub-agent/view"
	gomock "github.com/golang/mock/gomock"
)

// MockApihubClient is a mock of ApihubClient interface.
type MockApihubClient struct {
	ctrl     *gomock.Controller
	recorder *MockApihubClientMockRecorder
}

// MockApihubClientMockRecorder is the mock recorder for MockApihubClient.
type MockApihubClientMockRecorder struct {
	mock *MockApihubClient
}

// NewMockApihubClient creates a new mock instance.
func NewMockApihubClient(ctrl *gomock.Controller) *MockApihubClient {
	mock := &MockApihubClient{ctrl: ctrl}
	mock.recorder = &MockApihubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApihubClient) EXPECT() *MockApihubClientMockRecorder {
	return m.recorder
}

// CheckApiKeyValid mocks base method.
func (m *MockApihubClient) CheckApiKeyValid(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckApiKeyValid", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckApiKeyValid indicates an expected call of CheckApiKeyValid.
func (mr *MockApihubClientMockRecorder) CheckApiKeyValid(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckApiKeyValid", reflect.TypeOf((*MockApihubClient)(nil).CheckApiKeyValid), arg0)
}

// GetPackageByServiceName mocks base method.
func (m *MockApihubClient) GetPackageByServiceName(arg0 secctx.SecurityContext, arg1 string) (*view.SimplePackage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPackageByServiceName", arg0, arg1)
	ret0, _ := ret[0].(*view.SimplePackage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPackageByServiceName indicates an expected call of GetPackageByServiceName.
func (mr *MockApihubClientMockRecorder) GetPackageByServiceName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPackageByServiceName", reflect.TypeOf((*MockApihubClient)(nil).GetPackageByServiceName), arg0, arg1)
}

// GetRsaPublicKey mocks base method.
func (m *MockApihubClient) GetRsaPublicKey(arg0 secctx.SecurityContext) (*view.PublicKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRsaPublicKey", arg0)
	ret0, _ := ret[0].(*view.PublicKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRsaPublicKey indicates an expected call of GetRsaPublicKey.
func (mr *MockApihubClientMockRecorder) GetRsaPublicKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRsaPublicKey", reflect.TypeOf((*MockApihubClient)(nil).GetRsaPublicKey), arg0)
}

// GetUserPackagesPromoteStatuses mocks base method.
func (m *MockApihubClient) GetUserPackagesPromoteStatuses(arg0 secctx.SecurityContext, arg1 view.PackagesReq) (view.AvailablePackagePromoteStatuses, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPackagesPromoteStatuses", arg0, arg1)
	ret0, _ := ret[0].(view.AvailablePackagePromoteStatuses)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPackagesPromoteStatuses indicates an expected call of GetUserPackagesPromoteStatuses.
func (mr *MockApihubClientMockRecorder) GetUserPackagesPromoteStatuses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPackagesPromoteStatuses", reflect.TypeOf((*MockApihubClient)(nil).GetUserPackagesPromoteStatuses), arg0, arg1)
}

// GetVersions mocks base method.
func (m *MockApihubClient) GetVersions(arg0 secctx.SecurityContext, arg1 string, arg2, arg3 int) (*view.PublishedVersionsView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersions", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*view.PublishedVersionsView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersions indicates an expected call of GetVersions.
func (mr *MockApihubClientMockRecorder) GetVersions(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersions", reflect.TypeOf((*MockApihubClient)(nil).GetVersions), arg0, arg1, arg2, arg3)
}

// SendKeepaliveMessage mocks base method.
func (m *MockApihubClient) SendKeepaliveMessage(arg0 view.AgentKeepaliveMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendKeepaliveMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendKeepaliveMessage indicates an expected call of SendKeepaliveMessage.
func (mr *MockApihubClientMockRecorder) SendKeepaliveMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendKeepaliveMessage", reflect.TypeOf((*MockApihubClient)(nil).SendKeepaliveMessage), arg0)
}
