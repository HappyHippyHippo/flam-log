package log

import (
	"fmt"
	"strings"
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type stringSerializer struct{}

func newStringSerializer() Serializer {
	return &stringSerializer{}
}

func (stringSerializer) Close() error {
	return nil
}

func (stringSerializer) Serialize(
	timestamp time.Time,
	level Level,
	message string,
	_ flam.Bag,
) string {
	return fmt.Sprintf(
		"%s [%s] %s\n",
		timestamp.Format("2006-01-02T15:04:05.000-0700"),
		strings.ToUpper(LevelName[level]),
		message)
}
