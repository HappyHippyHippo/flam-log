package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	log "github.com/happyhippyhippo/flam-log"
)

// Serializer is a mock of LogSerializer interface.
type Serializer struct {
	ctrl     *gomock.Controller
	recorder *SerializerRecorder
}

// SerializerRecorder is the mock recorder for Serializer.
type SerializerRecorder struct {
	mock *Serializer
}

// NewSerializer creates a new mock instance.
func NewSerializer(ctrl *gomock.Controller) *Serializer {
	mock := &Serializer{ctrl: ctrl}
	mock.recorder = &SerializerRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Serializer) EXPECT() *SerializerRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Serializer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *SerializerRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Serializer)(nil).Close))
}

// Serialize mocks base method.
func (m *Serializer) Serialize(timestamp time.Time, level log.Level, message string, ctx flam.Bag) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Serialize", timestamp, level, message, ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// Serialize indicates an expected call of Serialize.
func (mr *SerializerRecorder) Serialize(timestamp, level, message, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Serialize", reflect.TypeOf((*Serializer)(nil).Serialize), timestamp, level, message, ctx)
}
