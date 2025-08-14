package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

type rotatingFileLogWriter struct {
	lock       sync.Locker
	disk       filesystem.Disk
	file       filesystem.File
	path       string
	timeFacade flamTime.Facade
	year       int
	month      time.Month
	day        int
	current    string
}

func newRotatingFileLogWriter(
	disk filesystem.Disk,
	path string,
	timeFacade flamTime.Facade,
) (io.Writer, error) {
	writer := &rotatingFileLogWriter{
		lock:       &sync.Mutex{},
		disk:       disk,
		path:       path,
		timeFacade: timeFacade,
	}

	if e := writer.rotate(timeFacade.Now()); e != nil {
		return nil, e
	}

	return writer, nil
}

func (writer *rotatingFileLogWriter) Write(
	output []byte,
) (int, error) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if e := writer.checkRotation(); e != nil {
		return 0, e
	}

	return writer.file.Write(output)
}

func (writer *rotatingFileLogWriter) Close() error {
	return writer.file.Close()
}

func (writer *rotatingFileLogWriter) checkRotation() error {
	now := writer.timeFacade.Now()
	if now.Day() != writer.day || now.Month() != writer.month || now.Year() != writer.year {
		return writer.rotate(now)
	}

	return nil
}

func (writer *rotatingFileLogWriter) rotate(
	now time.Time,
) error {
	writer.year = now.Year()
	writer.month = now.Month()
	writer.day = now.Day()
	writer.current = fmt.Sprintf(writer.path, now.Format("2006-01-02"))

	fp, e := writer.disk.OpenFile(writer.current, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if e != nil {
		return e
	}

	if writer.file != nil {
		_ = writer.file.Close()
	}
	writer.file = fp

	return nil
}
