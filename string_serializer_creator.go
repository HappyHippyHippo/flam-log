package log

import (
	flam "github.com/happyhippyhippo/flam"
)

type stringSerializerCreator struct{}

func newStringSerializerCreator() SerializerCreator {
	return &stringSerializerCreator{}
}

func (stringSerializerCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == SerializerDriverString
}

func (stringSerializerCreator) Create(
	_ flam.Bag,
) (Serializer, error) {
	return newStringSerializer(), nil
}
