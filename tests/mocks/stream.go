package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	log "github.com/happyhippyhippo/flam-log"
)

// Stream is a mock of LogStream interface.
type Stream struct {
	ctrl     *gomock.Controller
	recorder *StreamRecorder
}

// StreamRecorder is the mock recorder for Stream.
type StreamRecorder struct {
	mock *Stream
}

// NewStream creates a new mock instance.
func NewStream(ctrl *gomock.Controller) *Stream {
	mock := &Stream{ctrl: ctrl}
	mock.recorder = &StreamRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Stream) EXPECT() *StreamRecorder {
	return m.recorder
}

// AddChannel mocks base method.
func (m *Stream) AddChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddChannel indicates an expected call of AddChannel.
func (mr *StreamRecorder) AddChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChannel", reflect.TypeOf((*Stream)(nil).AddChannel), channel)
}

// Broadcast mocks base method.
func (m *Stream) Broadcast(timestamp time.Time, level log.Level, message string, ctx flam.Bag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Broadcast", timestamp, level, message, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Broadcast indicates an expected call of Broadcast.
func (mr *StreamRecorder) Broadcast(timestamp, level, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*Stream)(nil).Broadcast), timestamp, level, message, ctx)
}

// Close mocks base method.
func (m *Stream) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *StreamRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Stream)(nil).Close))
}

// GetLevel mocks base method.
func (m *Stream) GetLevel() log.Level {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLevel")
	ret0, _ := ret[0].(log.Level)
	return ret0
}

// GetLevel indicates an expected call of GetLevel.
func (mr *StreamRecorder) GetLevel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLevel", reflect.TypeOf((*Stream)(nil).GetLevel))
}

// HasChannel mocks base method.
func (m *Stream) HasChannel(channel string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasChannel", channel)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasChannel indicates an expected call of HasChannel.
func (mr *StreamRecorder) HasChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasChannel", reflect.TypeOf((*Stream)(nil).HasChannel), channel)
}

// ListChannels mocks base method.
func (m *Stream) ListChannels() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChannels")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListChannels indicates an expected call of ListChannels.
func (mr *StreamRecorder) ListChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChannels", reflect.TypeOf((*Stream)(nil).ListChannels))
}

// RemoveAllChannels mocks base method.
func (m *Stream) RemoveAllChannels() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAllChannels")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAllChannels indicates an expected call of RemoveAllChannels.
func (mr *StreamRecorder) RemoveAllChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllChannels", reflect.TypeOf((*Stream)(nil).RemoveAllChannels))
}

// RemoveChannel mocks base method.
func (m *Stream) RemoveChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveChannel indicates an expected call of RemoveChannel.
func (mr *StreamRecorder) RemoveChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChannel", reflect.TypeOf((*Stream)(nil).RemoveChannel), channel)
}

// SetLevel mocks base method.
func (m *Stream) SetLevel(level log.Level) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLevel", level)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLevel indicates an expected call of SetLevel.
func (mr *StreamRecorder) SetLevel(level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLevel", reflect.TypeOf((*Stream)(nil).SetLevel), level)
}

// Signal mocks base method.
func (m *Stream) Signal(timestamp time.Time, level log.Level, channel, message string, ctx flam.Bag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signal", timestamp, level, channel, message, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signal indicates an expected call of Signal.
func (mr *StreamRecorder) Signal(timestamp, level, channel, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*Stream)(nil).Signal), timestamp, level, channel, message, ctx)
}
