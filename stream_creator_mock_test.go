package log

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

type StreamCreatorMock struct {
	ctrl     *gomock.Controller
	recorder *StreamCreatorMockRecorder
}

type StreamCreatorMockRecorder struct {
	mock *StreamCreatorMock
}

func NewStreamCreatorMock(ctrl *gomock.Controller) *StreamCreatorMock {
	mock := &StreamCreatorMock{ctrl: ctrl}
	mock.recorder = &StreamCreatorMockRecorder{mock}
	return mock
}

func (m *StreamCreatorMock) EXPECT() *StreamCreatorMockRecorder {
	return m.recorder
}

func (m *StreamCreatorMock) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *StreamCreatorMockRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*StreamCreatorMock)(nil).Accept), config)
}

func (m *StreamCreatorMock) Create(config flam.Bag) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *StreamCreatorMockRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*StreamCreatorMock)(nil).Create), config)
}
