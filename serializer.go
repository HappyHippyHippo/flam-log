package log

import (
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type Serializer interface {
	Close() error

	Serialize(timestamp time.Time, level Level, message string, ctx flam.Bag) string
}
