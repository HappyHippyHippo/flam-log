package log

import (
	"io"
	"slices"
	"sort"
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type Stream interface {
	Close() error

	GetLevel() Level
	SetLevel(level Level) error

	HasChannel(channel string) bool
	ListChannels() []string
	AddChannel(channel string) error
	RemoveChannel(channel string) error
	RemoveAllChannels() error

	Signal(timestamp time.Time, level Level, channel, message string, ctx flam.Bag) error
	Broadcast(timestamp time.Time, level Level, message string, ctx flam.Bag) error
}

type stream struct {
	level      Level
	channels   []string
	serializer Serializer
	writer     io.Writer
	doClose    bool
}

func newStream(
	level Level,
	channels []string,
	serializer Serializer,
	writer io.Writer,
	doClose bool,
) *stream {
	return &stream{
		level:      level,
		channels:   channels,
		serializer: serializer,
		writer:     writer,
		doClose:    doClose,
	}
}

func (stream *stream) Close() error {
	if closer, ok := stream.writer.(io.Closer); stream.doClose && ok {
		return closer.Close()
	}

	return nil
}

func (stream *stream) GetLevel() Level {
	return stream.level
}

func (stream *stream) SetLevel(
	level Level,
) error {
	stream.level = level

	return nil
}

func (stream *stream) HasChannel(
	channel string,
) bool {
	return slices.Contains(stream.channels, channel)
}

func (stream *stream) ListChannels() []string {
	return stream.channels
}

func (stream *stream) AddChannel(
	channel string,
) error {
	if !stream.HasChannel(channel) {
		stream.channels = append(stream.channels, channel)
		sort.Strings(stream.channels)
	}

	return nil
}

func (stream *stream) RemoveChannel(
	channel string,
) error {
	stream.channels = slices.DeleteFunc(stream.channels, func(c string) bool {
		return c == channel
	})

	return nil
}

func (stream *stream) RemoveAllChannels() error {
	stream.channels = []string{}

	return nil
}

func (stream *stream) Signal(
	timestamp time.Time,
	level Level,
	channel,
	message string,
	ctx flam.Bag,
) error {
	if !stream.acceptChannel(channel) {
		return nil
	}

	ctx["channel"] = channel

	return stream.Broadcast(timestamp, level, message, ctx)
}

func (stream *stream) Broadcast(
	timestamp time.Time,
	level Level,
	message string,
	ctx flam.Bag,
) error {
	if stream.level < level || stream.level == None {
		return nil
	}

	serialized := stream.serializer.Serialize(timestamp, level, message, ctx)
	_, e := stream.writer.Write([]byte(serialized))

	return e
}

func (stream *stream) acceptChannel(
	channel string,
) bool {
	i := sort.SearchStrings(stream.channels, "*")
	if i != len(stream.channels) && stream.channels[i] == "*" {
		return true
	}

	i = sort.SearchStrings(stream.channels, channel)

	return i != len(stream.channels) && stream.channels[i] == channel
}
