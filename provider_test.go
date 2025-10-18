package log

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	"github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.log.provider", NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(t, NewProvider().Register(nil), flam.ErrNilReference)
	})

	t.Run("should successfully provide Facade", func(t *testing.T) {
		container := dig.New()

		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NotNil(t, facade)
		}))
	})
}

func Test_Provider_Boot(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(nil),
			flam.ErrNilReference)
	})

	t.Run("should use default boot values when not provided", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, Info, DefaultLevel)
		assert.Equal(t, "", DefaultSerializer)
		assert.Equal(t, "", DefaultDisk)
	})

	t.Run("should use configured default boot values when provided", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathDefaultLevel, Fatal)
		_ = config.Defaults.Set(PathDefaultSerializer, "my_serializer")
		_ = config.Defaults.Set(PathDefaultDisk, "my_disk")
		defer func() {
			DefaultLevel = Info
			DefaultSerializer = ""
			DefaultDisk = ""
			config.Defaults = flam.Bag{}
		}()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.Equal(t, Fatal, DefaultLevel)
		assert.Equal(t, "my_disk", DefaultDisk)
		assert.Equal(t, "my_serializer", DefaultSerializer)
	})

	t.Run("should return stream instantiation error", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathBoot, true)
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":   "invalid",
				"priority": 123,
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrInvalidResourceConfig)
	})

	t.Run("should return stream storing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathBoot, true)
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": SerializerDriverString,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverConsole,
				"serializer": "string",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		stm := NewStreamMock(ctrl)
		require.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.AddStream("my_stream", stm))
		}))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			ErrDuplicateStream)
	})

	t.Run("should correctly load the streams", func(t *testing.T) {
		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathBoot, true)
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": SerializerDriverString,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverConsole,
				"serializer": "string",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.True(t, facade.HasStream("my_stream"))
		}))
	})

	t.Run("should register a log flusher frequency config observer", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.True(t, facade.HasObserver("flam.log", PathFlusherFrequency))
		}))
	})

	t.Run("should update the log flusher observer trigger when the config flusher frequency changes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathFlusherFrequency, time.Millisecond*10)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		p := NewProvider()
		require.NoError(t, p.Register(container))

		wg := &sync.WaitGroup{}
		stream := NewStreamMock(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), Info, "channel", "message", flam.Bag{"ctx": "value"}).
			DoAndReturn(func(
				timestamp time.Time,
				level Level,
				channel, message string,
				ctx flam.Bag,
			) error {
				wg.Done()
				return nil
			}).
			Times(2)
		require.NoError(t, container.Invoke(func(facade Facade) error {
			return facade.AddStream("my_stream", stream)
		}))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		wg.Add(1)
		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Signal(Info, "channel", "message", flam.Bag{"ctx": "value"}))
		}))
		wg.Wait() // wait for the log stream flush

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(PathFlusherFrequency, time.Millisecond*20))
		}))

		wg.Add(1)
		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Signal(Info, "channel", "message", flam.Bag{"ctx": "value"}))
		}))
		wg.Wait() // wait for the log stream flush after the frequency change
	})

	t.Run("should not recreate recurring trigger when config frequency change to a value that is not a duration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathFlusherFrequency, time.Millisecond*10)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		timeFacade := NewTimeFacadeMock(ctrl)
		timeFacade.EXPECT().
			NewRecurringTrigger(time.Millisecond*10, gomock.Any()).
			Return(NewTriggerMock(ctrl), nil).
			Times(1)
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))

		p := NewProvider()
		require.NoError(t, p.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(PathFlusherFrequency, "string"))
		}))
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error on nil container", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return serializer closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		expectedError := errors.New("close error")
		serializer := NewSerializerMock(ctrl)
		serializer.EXPECT().Close().Return(expectedError).Times(1)
		require.NoError(t, container.Invoke(func(facade Facade) error {
			return facade.AddSerializer("serializer", serializer)
		}))

		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			p.(flam.ClosableProvider).Close(container),
			expectedError)
	})

	t.Run("should close a loaded string serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": SerializerDriverString,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverConsole,
				"serializer": "string",
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})

	t.Run("should close a loaded json serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverConsole,
				"serializer": "json",
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})

	t.Run("should return stream closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(config.PathBoot, true)
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		expectedError := errors.New("close error")
		stm := NewStreamMock(ctrl)
		stm.EXPECT().Close().Return(expectedError).Times(1)

		stmCreatorConfig := flam.Bag{"id": "my_stream", "driver": "mock"}
		stmCreator := NewStreamCreatorMock(ctrl)
		stmCreator.EXPECT().Accept(stmCreatorConfig).Return(true).Times(1)
		stmCreator.EXPECT().Create(stmCreatorConfig).Return(stm, nil).Times(1)
		require.NoError(t, container.Provide(func() StreamCreator {
			return stmCreator
		}, dig.Group(StreamCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(t, p.(flam.ClosableProvider).Close(container), expectedError)
	})

	t.Run("should close the config observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathBoot, true)
		_ = config.Defaults.Set(PathFlusherFrequency, 1000)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		p := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, p.Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, p.(flam.BootableProvider).Boot(container))

		assert.NoError(t, p.(flam.ClosableProvider).Close(container))
	})
}
