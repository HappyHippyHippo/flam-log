package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	flam "github.com/happyhippyhippo/flam"
	log "github.com/happyhippyhippo/flam-log"
)

// SerializerCreator is a mock of LogSerializerCreator interface.
type SerializerCreator struct {
	ctrl     *gomock.Controller
	recorder *SerializerCreatorRecorder
}

// SerializerCreatorRecorder is the mock recorder for SerializerCreator.
type SerializerCreatorRecorder struct {
	mock *SerializerCreator
}

// NewSerializerCreator creates a new mock instance.
func NewSerializerCreator(ctrl *gomock.Controller) *SerializerCreator {
	mock := &SerializerCreator{ctrl: ctrl}
	mock.recorder = &SerializerCreatorRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SerializerCreator) EXPECT() *SerializerCreatorRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *SerializerCreator) Accept(config flam.Bag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", config)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *SerializerCreatorRecorder) Accept(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*SerializerCreator)(nil).Accept), config)
}

// Create mocks base method.
func (m *SerializerCreator) Create(config flam.Bag) (log.Serializer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", config)
	ret0, _ := ret[0].(log.Serializer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *SerializerCreatorRecorder) Create(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*SerializerCreator)(nil).Create), config)
}
