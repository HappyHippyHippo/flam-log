package log

import (
	"os"
	"sort"

	flam "github.com/happyhippyhippo/flam"
)

type consoleStreamCreator struct {
	streamCreator
}

func newConsoleStreamCreator(
	serializerFactory serializerFactory,
) StreamCreator {
	return &consoleStreamCreator{
		streamCreator: streamCreator{
			serializerFactory: serializerFactory,
		},
	}
}

func (consoleStreamCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == StreamDriverConsole
}

func (creator consoleStreamCreator) Create(
	config flam.Bag,
) (Stream, error) {
	serializerId := config.String("serializer", DefaultSerializer)
	serializer, e := creator.serializerFactory.Get(serializerId)
	if e != nil {
		return nil, e
	}

	channels := creator.getChannels(config.Slice("channels"))
	sort.Strings(channels)

	return newStream(
		LevelFrom(config.Get("level"), DefaultLevel),
		channels,
		serializer,
		os.Stdout,
		false), nil
}
