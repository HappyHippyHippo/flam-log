package tests

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	log "github.com/happyhippyhippo/flam-log"
	mocks "github.com/happyhippyhippo/flam-log/tests/mocks"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, log.NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.log.provider", log.NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(t, log.NewProvider().Register(nil), flam.ErrNilReference)
	})

	t.Run("should successfully provide Facade", func(t *testing.T) {
		container := dig.New()

		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NotNil(t, facade)
		}))
	})
}

func Test_Provider_Boot(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			log.NewProvider().(flam.BootableProvider).Boot(nil),
			flam.ErrNilReference)
	})

	t.Run("should use default boot values when not provided", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, log.Info, log.DefaultLevel)
		assert.Equal(t, "", log.DefaultSerializer)
		assert.Equal(t, "", log.DefaultDisk)
	})

	t.Run("should use configured default boot values when provided", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathDefaultLevel, log.Fatal)
		_ = config.Defaults.Set(log.PathDefaultSerializer, "my_serializer")
		_ = config.Defaults.Set(log.PathDefaultDisk, "my_disk")
		defer func() {
			log.DefaultLevel = log.Info
			log.DefaultSerializer = ""
			log.DefaultDisk = ""
			config.Defaults = flam.Bag{}
		}()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, log.Fatal, log.DefaultLevel)
		assert.Equal(t, "my_disk", log.DefaultDisk)
		assert.Equal(t, "my_serializer", log.DefaultSerializer)
	})

	t.Run("should return stream instantiation error", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathBoot, true)
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":   "invalid",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			log.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return stream storing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathBoot, true)
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		require.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.AddStream("my_stream", stream))
		}))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			log.NewProvider().(flam.BootableProvider).Boot(container),
			log.ErrDuplicateStream)
	})

	t.Run("should correctly load the streams", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathBoot, true)
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.True(t, facade.HasStream("my_stream"))
		}))
	})
}

func Test_Provider_Run(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			log.NewProvider().(flam.RunnableProvider).Run(nil),
			flam.ErrNilReference)
	})

	t.Run("should register a log flusher frequency config observer", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		require.NoError(t, log.NewProvider().(flam.RunnableProvider).Run(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.True(t, facade.HasObserver("flam.log", log.PathFlusherFrequency))
		}))
	})

	t.Run("should update the log flusher observer trigger when the config flusher frequency changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathFlusherFrequency, time.Millisecond*10)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		provider := log.NewProvider()
		require.NoError(t, provider.Register(container))

		wg := &sync.WaitGroup{}
		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{"ctx": "value"}).
			DoAndReturn(func(
				timestamp time.Time,
				level log.Level,
				channel, message string,
				ctx flam.Bag,
			) error {
				wg.Done()
				return nil
			}).
			Times(2)
		require.NoError(t, container.Invoke(func(facade log.Facade) error {
			return facade.AddStream("my_stream", stream)
		}))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		wg.Add(1)
		require.NoError(t, provider.(flam.RunnableProvider).Run(container))
		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.Signal(log.Info, "channel", "message", flam.Bag{"ctx": "value"}))
		}))
		wg.Wait() // wait for the log stream flush

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(log.PathFlusherFrequency, time.Millisecond*20))
		}))

		wg.Add(1)
		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.Signal(log.Info, "channel", "message", flam.Bag{"ctx": "value"}))
		}))
		wg.Wait() // wait for the log stream flush after the frequency change
	})

	t.Run("should not recreate recurring trigger when config frequency change to a value that is not a duration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathFlusherFrequency, time.Millisecond*10)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		timeFacade := mocks.NewTimeFacade(ctrl)
		timeFacade.EXPECT().
			NewRecurringTrigger(time.Millisecond*10, gomock.Any()).
			Return(mocks.NewTrigger(ctrl), nil).
			Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		provider := log.NewProvider()
		require.NoError(t, provider.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(log.PathFlusherFrequency, "string"))
		}))
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			log.NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return serializer closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		provider := log.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		expectedError := errors.New("close error")
		serializer := mocks.NewSerializer(ctrl)
		serializer.EXPECT().Close().Return(expectedError).Times(1)
		require.NoError(t, container.Invoke(func(facade log.Facade) error {
			return facade.AddSerializer("serializer", serializer)
		}))

		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			provider.(flam.ClosableProvider).Close(container),
			expectedError)
	})

	t.Run("should close a loaded string serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := log.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})

	t.Run("should close a loaded json serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": log.SerializerDriverJson,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "json",
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := log.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})

	t.Run("should return stream closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := log.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		expectedError := errors.New("close error")
		stream := mocks.NewStream(ctrl)
		stream.EXPECT().Close().Return(expectedError).Times(1)

		streamCreatorConfig := flam.Bag{"id": "my_stream", "driver": "mock"}
		streamCreator := mocks.NewStreamCreator(ctrl)
		streamCreator.EXPECT().Accept(streamCreatorConfig).Return(true).Times(1)
		streamCreator.EXPECT().Create(streamCreatorConfig).Return(stream, nil).Times(1)
		require.NoError(t, container.Provide(func() log.StreamCreator {
			return streamCreator
		}, dig.Group(log.StreamCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.ErrorIs(t, provider.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close the config observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathBoot, true)
		_ = config.Defaults.Set(log.PathFlusherFrequency, 1000)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := log.NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		require.NoError(t, provider.(flam.RunnableProvider).Run(container))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})
}
