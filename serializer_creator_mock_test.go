package log

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

type SerializerCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *SerializerCreatorMockRecorder
}

type SerializerCreatorMockRecorder struct {
	mock *SerializerCreatorMock
}

func NewSerializerCreatorMock(ctrl *gomock.Controller) *SerializerCreatorMock {
	mock := &SerializerCreatorMock{ctrl: ctrl}
	mock.recorder = &SerializerCreatorMockRecorder{mock}
	return mock
}

func (m *SerializerCreatorMock) EXPECT() *SerializerCreatorMockRecorder {
	return m.recorder
}

func (m *SerializerCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *SerializerCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*SerializerCreatorMock)(nil).Accept), config)
}

func (m *SerializerCreatorMock) Create(config flam.Bag) (Serializer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(Serializer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *SerializerCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*SerializerCreatorMock)(nil).Create), config)
}
