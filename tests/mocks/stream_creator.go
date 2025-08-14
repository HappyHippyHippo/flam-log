package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	log "github.com/happyhippyhippo/flam-log"
)

// StreamCreator is a mock of LogStreamCreator interface.
type StreamCreator struct {
	ctrl     *gomock.Controller
	recorder *StreamCreatorRecorder
}

// StreamCreatorRecorder is the mock recorder for StreamCreator.
type StreamCreatorRecorder struct {
	mock *StreamCreator
}

// NewStreamCreator creates a new mock instance.
func NewStreamCreator(ctrl *gomock.Controller) *StreamCreator {
	mock := &StreamCreator{ctrl: ctrl}
	mock.recorder = &StreamCreatorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *StreamCreator) EXPECT() *StreamCreatorRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *StreamCreator) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *StreamCreatorRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*StreamCreator)(nil).Accept), config)
}

// Create mocks base method.
func (m *StreamCreator) Create(config flam.Bag) (log.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(log.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *StreamCreatorRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*StreamCreator)(nil).Create), config)
}
