package log

import (
	flam "github.com/happyhippyhippo/flam"
)

type jsonSerializerCreator struct{}

func newJsonSerializerCreator() SerializerCreator {
	return &jsonSerializerCreator{}
}

func (jsonSerializerCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SerializerDriverJson
}

func (jsonSerializerCreator) Create(
	_ flam.Bag,
) (Serializer, error) {
	return newJsonSerializer(), nil
}
