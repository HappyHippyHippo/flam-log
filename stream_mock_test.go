package log

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
)

type StreamMock struct {
	ctrl     *gomock.Controller
	recorder *StreamMockRecorder
}

type StreamMockRecorder struct {
	mock *StreamMock
}

func NewStreamMock(ctrl *gomock.Controller) *StreamMock {
	mock := &StreamMock{ctrl: ctrl}
	mock.recorder = &StreamMockRecorder{mock}
	return mock
}

func (m *StreamMock) EXPECT() *StreamMockRecorder {
	return m.recorder
}

func (m *StreamMock) AddChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) AddChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChannel", reflect.TypeOf((*StreamMock)(nil).AddChannel), channel)
}

func (m *StreamMock) Broadcast(timestamp time.Time, level Level, message string, ctx flam.Bag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Broadcast", timestamp, level, message, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) Broadcast(timestamp, level, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*StreamMock)(nil).Broadcast), timestamp, level, message, ctx)
}

func (m *StreamMock) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*StreamMock)(nil).Close))
}

func (m *StreamMock) GetLevel() Level {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLevel")
	ret0, _ := ret[0].(Level)
	return ret0
}

func (mr *StreamMockRecorder) GetLevel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLevel", reflect.TypeOf((*StreamMock)(nil).GetLevel))
}

func (m *StreamMock) HasChannel(channel string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasChannel", channel)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (mr *StreamMockRecorder) HasChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasChannel", reflect.TypeOf((*StreamMock)(nil).HasChannel), channel)
}

func (m *StreamMock) ListChannels() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChannels")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (mr *StreamMockRecorder) ListChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChannels", reflect.TypeOf((*StreamMock)(nil).ListChannels))
}

func (m *StreamMock) RemoveAllChannels() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAllChannels")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) RemoveAllChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllChannels", reflect.TypeOf((*StreamMock)(nil).RemoveAllChannels))
}

func (m *StreamMock) RemoveChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) RemoveChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChannel", reflect.TypeOf((*StreamMock)(nil).RemoveChannel), channel)
}

func (m *StreamMock) SetLevel(level Level) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLevel", level)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) SetLevel(level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLevel", reflect.TypeOf((*StreamMock)(nil).SetLevel), level)
}

func (m *StreamMock) Signal(timestamp time.Time, level Level, channel, message string, ctx flam.Bag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signal", timestamp, level, channel, message, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *StreamMockRecorder) Signal(timestamp, level, channel, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*StreamMock)(nil).Signal), timestamp, level, channel, message, ctx)
}
