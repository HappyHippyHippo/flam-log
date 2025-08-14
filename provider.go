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

	var e error
	provide := func(constructor any, opts ...dig.ProvideOption) bool {
		e = container.Provide(constructor, opts...)
		return e == nil
	}

	_ = provide(newStringSerializerCreator, dig.Group(SerializerCreatorGroup)) &&
		provide(newJsonSerializerCreator, dig.Group(SerializerCreatorGroup)) &&
		provide(newSerializerFactory) &&
		provide(newConsoleStreamCreator, dig.Group(StreamCreatorGroup)) &&
		provide(newFileStreamCreator, dig.Group(StreamCreatorGroup)) &&
		provide(newRotatingFileStreamCreator, dig.Group(StreamCreatorGroup)) &&
		provide(newStreamFactory) &&
		provide(newManager) &&
		provide(newFacade)

	return e
}

func (*provider) Boot(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		configFacade config.Facade,
		manager *manager,
		streamFactory steamFactory,
	) error {
		DefaultLevel = LevelFrom(configFacade.Get(PathDefaultLevel), DefaultLevel)
		DefaultSerializer = configFacade.String(PathDefaultSerializer, DefaultSerializer)
		DefaultDisk = configFacade.String(PathDefaultDisk, DefaultDisk)

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
	})
}

func (provider *provider) Run(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		timeFacade flamTime.Facade,
		configFacade config.Facade,
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
	})
}

func (provider *provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		streamFactory steamFactory,
		serializerFactory serializerFactory,
	) error {
		if provider.flusher != nil {
			_ = provider.flusher.Close()
		}

		if e := streamFactory.Close(); e != nil {
			return e
		}

		if e := serializerFactory.Close(); e != nil {
			return e
		}

		return nil
	})
}
