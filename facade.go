package log

import (
	flam "github.com/happyhippyhippo/flam"
)

type Facade interface {
	Signal(level Level, channel, message string, ctx ...flam.Bag) error
	Broadcast(level Level, message string, ctx ...flam.Bag) error
	FatalSignal(channel, message string, ctx ...flam.Bag) error
	FatalBroadcast(message string, ctx ...flam.Bag) error
	ErrorSignal(channel, message string, ctx ...flam.Bag) error
	ErrorBroadcast(message string, ctx ...flam.Bag) error
	WarningSignal(channel, message string, ctx ...flam.Bag) error
	WarningBroadcast(message string, ctx ...flam.Bag) error
	NoticeSignal(channel, message string, ctx ...flam.Bag) error
	NoticeBroadcast(message string, ctx ...flam.Bag) error
	InfoSignal(channel, message string, ctx ...flam.Bag) error
	InfoBroadcast(message string, ctx ...flam.Bag) error
	DebugSignal(channel, message string, ctx ...flam.Bag) error
	DebugBroadcast(message string, ctx ...flam.Bag) error
	Flush() error

	HasSerializer(id string) bool
	ListSerializers() []string
	GetSerializer(id string) (Serializer, error)
	AddSerializer(id string, serializer Serializer) error

	HasStream(id string) bool
	ListStreams() []string
	GetStream(id string) (Stream, error)
	AddStream(id string, stream Stream) error
	RemoveStream(id string) error
	RemoveAllStreams() error
}

type facade struct {
	serializerFactory serializerFactory
	manager           *manager
}

func newFacade(
	serializerFactory serializerFactory,
	manager *manager,
) Facade {
	return &facade{
		serializerFactory: serializerFactory,
		manager:           manager,
	}
}

func (facade *facade) Signal(
	level Level,
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(level, channel, message, ctx...)
}

func (facade *facade) Broadcast(
	level Level,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(level, message, ctx...)
}

func (facade *facade) FatalSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Fatal, channel, message, ctx...)
}

func (facade *facade) FatalBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Fatal, message, ctx...)
}

func (facade *facade) ErrorSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Error, channel, message, ctx...)
}

func (facade *facade) ErrorBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Error, message, ctx...)
}

func (facade *facade) WarningSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Warning, channel, message, ctx...)
}

func (facade *facade) WarningBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Warning, message, ctx...)
}

func (facade *facade) NoticeSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Notice, channel, message, ctx...)
}

func (facade *facade) NoticeBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Notice, message, ctx...)
}

func (facade *facade) InfoSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Info, channel, message, ctx...)
}

func (facade *facade) InfoBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Info, message, ctx...)
}

func (facade *facade) DebugSignal(
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Signal(Debug, channel, message, ctx...)
}

func (facade *facade) DebugBroadcast(
	message string,
	ctx ...flam.Bag,
) error {
	return facade.manager.Broadcast(Debug, message, ctx...)
}

func (facade *facade) Flush() error {
	return facade.manager.Flush()
}

func (facade *facade) HasSerializer(
	id string,
) bool {
	return facade.serializerFactory.Has(id)
}

func (facade *facade) ListSerializers() []string {
	return facade.serializerFactory.List()
}

func (facade *facade) GetSerializer(
	id string,
) (Serializer, error) {
	return facade.serializerFactory.Get(id)
}

func (facade *facade) AddSerializer(
	id string,
	serializer Serializer,
) error {
	return facade.serializerFactory.Add(id, serializer)
}

func (facade *facade) HasStream(
	id string,
) bool {
	return facade.manager.HasStream(id)
}

func (facade *facade) ListStreams() []string {
	return facade.manager.ListStreams()
}

func (facade *facade) GetStream(
	id string,
) (Stream, error) {
	return facade.manager.GetStream(id)
}

func (facade *facade) AddStream(
	id string,
	stream Stream,
) error {
	return facade.manager.AddStream(id, stream)
}

func (facade *facade) RemoveStream(
	id string,
) error {
	return facade.manager.RemoveStream(id)
}

func (facade *facade) RemoveAllStreams() error {
	return facade.manager.RemoveAllStreams()
}
