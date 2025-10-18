package log

import (
	"time"

	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	flamTime "github.com/happyhippyhippo/flam-time"
)

type provider struct {
	flusher flamTime.Trigger
}

func NewProvider() flam.Provider {
	return &provider{}
}

func (*provider) Id() string {
	return providerId
}

func (*provider) Register(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	registerer := flam.NewRegisterer()
	registerer.Queue(newStringSerializerCreator, dig.Group(SerializerCreatorGroup))
	registerer.Queue(newJsonSerializerCreator, dig.Group(SerializerCreatorGroup))
	registerer.Queue(newSerializerFactory)
	registerer.Queue(newConsoleStreamCreator, dig.Group(StreamCreatorGroup))
	registerer.Queue(newFileStreamCreator, dig.Group(StreamCreatorGroup))
	registerer.Queue(newRotatingFileStreamCreator, dig.Group(StreamCreatorGroup))
	registerer.Queue(newStreamFactory)
	registerer.Queue(newManager)
	registerer.Queue(newFacade)

	return registerer.Run(container)
}

func (provider *provider) Boot(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	executor := flam.NewExecutor()
	executor.Queue(provider.bootDefaults)
	executor.Queue(provider.bootStreams)
	executor.Queue(provider.bootFlusher)

	return executor.Run(container)
}

func (provider *provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	executor := flam.NewExecutor()
	executor.Queue(provider.closeFlusher)
	executor.Queue(provider.closeStreamFactory)
	executor.Queue(provider.closeSerializerFactory)

	return executor.Run(container)
}

func (*provider) bootDefaults(
	configFacade config.Facade,
) error {
	DefaultLevel = LevelFrom(configFacade.Get(PathDefaultLevel), DefaultLevel)
	DefaultSerializer = configFacade.String(PathDefaultSerializer, DefaultSerializer)
	DefaultDisk = configFacade.String(PathDefaultDisk, DefaultDisk)

	return nil
}

func (*provider) bootStreams(
	configFacade config.Facade,
	streamFactory steamFactory,
	manager *manager,
) error {
	if configFacade.Bool(PathBoot) {
		for id := range configFacade.Bag(PathStreams) {
			stream, e := streamFactory.Get(id)
			if e != nil {
				return e
			}

			if e = manager.AddStream(id, stream); e != nil {
				return e
			}
		}
	}

	return nil
}

func (provider *provider) bootFlusher(
	configFacade config.Facade,
	timeFacade flamTime.Facade,
	manager *manager,
) error {
	frequency := configFacade.Duration(PathFlusherFrequency)
	if frequency != time.Duration(0) {
		provider.flusher, _ = timeFacade.NewRecurringTrigger(frequency, func() error {
			return manager.Flush()
		})
	}

	return configFacade.AddObserver(
		"flam.log",
		PathFlusherFrequency,
		func(old, new any) {
			frequency, ok := new.(time.Duration)
			if !ok {
				return
			}

			_ = provider.flusher.Close()

			provider.flusher, _ = timeFacade.NewRecurringTrigger(frequency, func() error {
				return manager.Flush()
			})
		},
	)
}

func (provider *provider) closeFlusher() error {
	if provider.flusher == nil {
		return nil
	}

	return provider.flusher.Close()
}

func (*provider) closeStreamFactory(
	streamFactory steamFactory,
) error {
	return streamFactory.Close()
}

func (*provider) closeSerializerFactory(
	serializerFactory serializerFactory,
) error {
	return serializerFactory.Close()
}
