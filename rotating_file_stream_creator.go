package log

import (
	"sort"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	time "github.com/happyhippyhippo/flam-time"
)

type rotatingFileStreamCreator struct {
	fileStreamCreator

	timeFacade time.Facade
}

func newRotatingFileStreamCreator(
	timeFacade time.Facade,
	fileSystemFacade filesystem.Facade,
	serializerFactory serializerFactory,
) StreamCreator {
	return &rotatingFileStreamCreator{
		fileStreamCreator: fileStreamCreator{
			streamCreator: streamCreator{
				serializerFactory: serializerFactory,
			},
			fileSystemFacade: fileSystemFacade,
		},
		timeFacade: timeFacade,
	}
}

func (creator rotatingFileStreamCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == StreamDriverRotatingFile &&
		config.Has("path")
}

func (creator rotatingFileStreamCreator) Create(
	config flam.Bag,
) (Stream, error) {
	serializerId := config.String("serializer", DefaultSerializer)
	serializer, e := creator.serializerFactory.Get(serializerId)
	if e != nil {
		return nil, e
	}

	diskId := config.String("disk", DefaultDisk)
	disk, e := creator.fileSystemFacade.GetDisk(diskId)
	if e != nil {
		return nil, e
	}

	file, e := newRotatingFileLogWriter(disk, config.String("path"), creator.timeFacade)
	if e != nil {
		return nil, e
	}

	channels := creator.getChannels(config.Slice("channels"))
	sort.Strings(channels)

	return newStream(
		LevelFrom(config.Get("level"), DefaultLevel),
		channels,
		serializer,
		file,
		true), nil
}
