package log

import (
	flam "github.com/happyhippyhippo/flam"
)

type SerializerCreator interface {
	Accept(config flam.Bag) bool
	Create(config flam.Bag) (Serializer, error)
}
