package log

import (
	"sync"
	"time"

	flam "github.com/happyhippyhippo/flam"
)

type regEntry struct {
	timestamp time.Time
	level     Level
	channel   string
	message   string
	ctx       flam.Bag
}

type manager struct {
	streams map[string]Stream
	buffer  []regEntry
	mutex   sync.Locker
}

func newManager() *manager {
	return &manager{
		streams: map[string]Stream{},
		buffer:  []regEntry{},
		mutex:   &sync.Mutex{},
	}
}

func (manager *manager) Signal(
	level Level,
	channel,
	message string,
	ctx ...flam.Bag,
) error {
	context := flam.Bag{}
	for _, c := range ctx {
		context.Merge(c)
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.buffer = append(manager.buffer, regEntry{
		timestamp: time.Now(),
		level:     level,
		channel:   channel,
		message:   message,
		ctx:       context})

	return nil
}

func (manager *manager) Broadcast(
	level Level,
	message string,
	ctx ...flam.Bag,
) error {
	context := flam.Bag{}
	for _, c := range ctx {
		context.Merge(c)
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.buffer = append(manager.buffer, regEntry{
		timestamp: time.Now(),
		level:     level,
		channel:   "",
		message:   message,
		ctx:       context,
	})

	return nil
}

func (manager *manager) Flush() error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for _, entry := range manager.buffer {
		for _, stream := range manager.streams {
			if entry.channel != "" {
				e := stream.Signal(
					entry.timestamp,
					entry.level,
					entry.channel,
					entry.message,
					entry.ctx)
				if e != nil {
					return e
				}
			} else {
				e := stream.Broadcast(
					entry.timestamp,
					entry.level,
					entry.message,
					entry.ctx)
				if e != nil {
					return e
				}
			}
		}
	}

	manager.buffer = []regEntry{}

	return nil
}

func (manager *manager) HasStream(
	id string,
) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	_, ok := manager.streams[id]

	return ok
}

func (manager *manager) ListStreams() []string {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	var list []string
	for id := range manager.streams {
		list = append(list, id)
	}

	return list
}

func (manager *manager) GetStream(
	id string,
) (Stream, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	stream, ok := manager.streams[id]
	if !ok {
		return nil, newErrStreamNotFound(id)
	}

	return stream, nil
}

func (manager *manager) AddStream(
	id string,
	stream Stream,
) error {
	switch {
	case stream == nil:
		return newErrNilReference("stream")
	case manager.HasStream(id):
		return newErrDuplicateStream(id)
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.streams[id] = stream

	return nil
}

func (manager *manager) RemoveStream(
	id string,
) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if _, ok := manager.streams[id]; !ok {
		return newErrStreamNotFound(id)
	}

	if e := manager.streams[id].Close(); e != nil {
		return e
	}
	delete(manager.streams, id)

	return nil
}

func (manager *manager) RemoveAllStreams() error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for _, stream := range manager.streams {
		if e := stream.Close(); e != nil {
			return e
		}
	}

	manager.streams = map[string]Stream{}

	return nil
}
