package log

import (
	"encoding/json"
	"strings"
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type jsonSerializer struct{}

func newJsonSerializer() Serializer {
	return &jsonSerializer{}
}

func (jsonSerializer) Close() error {
	return nil
}

func (jsonSerializer) Serialize(
	timestamp time.Time,
	level Level,
	message string,
	ctx flam.Bag,
) string {
	ctx["time"] = timestamp.Format("2006-01-02T15:04:05.000-0700")
	ctx["level"] = strings.ToUpper(LevelName[level])
	ctx["message"] = message
	bytes, _ := json.Marshal(ctx)
	bytes = append(bytes, '\n')

	return string(bytes)
}
