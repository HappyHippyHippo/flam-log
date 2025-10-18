package log

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

type SerializerMock struct {
	ctrl     *gomock.Controller
	recorder *SerializerMockRecorder
}

type SerializerMockRecorder struct {
	mock *SerializerMock
}

func NewSerializerMock(ctrl *gomock.Controller) *SerializerMock {
	mock := &SerializerMock{ctrl: ctrl}
	mock.recorder = &SerializerMockRecorder{mock}
	return mock
}

func (m *SerializerMock) EXPECT() *SerializerMockRecorder {
	return m.recorder
}

func (m *SerializerMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *SerializerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*SerializerMock)(nil).Close))
}

func (m *SerializerMock) Serialize(timestamp time.Time, level Level, message string, ctx flam.Bag) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serialize", timestamp, level, message, ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

func (mr *SerializerMockRecorder) Serialize(timestamp, level, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serialize", reflect.TypeOf((*SerializerMock)(nil).Serialize), timestamp, level, message, ctx)
}
