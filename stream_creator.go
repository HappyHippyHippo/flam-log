package log

import (
	flam "github.com/happyhippyhippo/flam"
)

type StreamCreator interface {
	Accept(config flam.Bag) bool
	Create(config flam.Bag) (Stream, error)
}

type streamCreator struct {
	serializerFactory flam.Factory[Serializer]
}

func (streamCreator) getChannels(
	list []any,
) []string {
	var result []string
	for _, channel := range list {
		if typedChannel, ok := channel.(string); ok {
			result = append(result, typedChannel)
		}
	}

	return result
}
