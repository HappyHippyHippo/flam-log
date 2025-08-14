package log

import (
	"os"
	"sort"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
)

type fileStreamCreator struct {
	streamCreator

	fileSystemFacade filesystem.Facade
}

func newFileStreamCreator(
	fileSystemFacade filesystem.Facade,
	serializerFactory serializerFactory,
) StreamCreator {
	return &fileStreamCreator{
		streamCreator: streamCreator{
			serializerFactory: serializerFactory,
		},
		fileSystemFacade: fileSystemFacade,
	}
}

func (fileStreamCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == StreamDriverFile &&
		config.Has("path")
}

func (creator fileStreamCreator) Create(
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

	file, e := disk.OpenFile(config.String("path"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
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
