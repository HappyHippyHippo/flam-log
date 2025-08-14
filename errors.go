package log

import (
	"errors"

	flam "github.com/happyhippyhippo/flam"
)

var (
	ErrStreamNotFound  = errors.New("log stream not found")
	ErrDuplicateStream = errors.New("duplicate log stream")
)

func newErrNilReference(
	field string,
) error {
	return flam.NewErrorFrom(
		flam.ErrNilReference,
		field)
}

func newErrStreamNotFound(
	id string,
) error {
	return flam.NewErrorFrom(
		ErrStreamNotFound,
		id)
}

func newErrDuplicateStream(
	id string,
) error {
	return flam.NewErrorFrom(
		ErrDuplicateStream,
		id)
}
