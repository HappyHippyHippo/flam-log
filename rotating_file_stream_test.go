package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	flamTime "github.com/happyhippyhippo/flam-time"
)

func Test_rotatingFileStream(t *testing.T) {
	t.Run("should ignore config without path field", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "invalid",
			}})
		_ = config.Defaults.Set(PathBoot, true)
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

	t.Run("should return serialization creation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "invalid",
				"path":       "/path",
			}})
		_ = config.Defaults.Set(PathBoot, true)
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
			flam.ErrUnknownResource)
	})

	t.Run("should return disk creation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"disk":       "invalid",
				"path":       "/path",
			}})
		_ = config.Defaults.Set(PathBoot, true)
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
			flam.ErrUnknownResource)
	})

	t.Run("should return file opening error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"disk":       "mock",
				"path":       "/file-%s",
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		expectedError := fmt.Errorf("error")
		now := time.Now()
		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile(fmt.Sprintf("/file-%s", now.Format("2006-01-02")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(nil, expectedError)

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			NewProvider().(flam.BootableProvider).Boot(container),
			expectedError)
	})

	t.Run("should correctly handle the stream level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "debug",
				"disk":       "mock",
				"path":       "/file-%s",
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			stream, e := facade.GetStream("my_stream")
			assert.NotNil(t, stream)
			assert.NoError(t, e)

			assert.Equal(t, Debug, stream.GetLevel())

			assert.NoError(t, stream.SetLevel(Info))
			assert.Equal(t, Info, stream.GetLevel())
		}))
	})

	t.Run("should correctly handle the channel list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "debug",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_2", "channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			stream, e := facade.GetStream("my_stream")
			require.NotNil(t, stream)
			require.NoError(t, e)

			assert.True(t, stream.HasChannel("channel_1"))
			assert.True(t, stream.HasChannel("channel_2"))
			assert.False(t, stream.HasChannel("channel_3"))
			assert.ElementsMatch(t, stream.ListChannels(), []string{"channel_1", "channel_2"})

			require.NoError(t, stream.AddChannel("channel_3"))
			assert.True(t, stream.HasChannel("channel_1"))
			assert.True(t, stream.HasChannel("channel_2"))
			assert.True(t, stream.HasChannel("channel_3"))
			assert.ElementsMatch(t, stream.ListChannels(), []string{"channel_1", "channel_2", "channel_3"})

			require.NoError(t, stream.RemoveChannel("channel_2"))
			assert.True(t, stream.HasChannel("channel_1"))
			assert.False(t, stream.HasChannel("channel_2"))
			assert.True(t, stream.HasChannel("channel_3"))
			assert.ElementsMatch(t, stream.ListChannels(), []string{"channel_1", "channel_3"})

			require.NoError(t, stream.RemoveAllChannels())
			assert.False(t, stream.HasChannel("channel_1"))
			assert.False(t, stream.HasChannel("channel_2"))
			assert.False(t, stream.HasChannel("channel_3"))
			assert.ElementsMatch(t, stream.ListChannels(), []string{})
		}))
	})

	t.Run("should correctly handle the stream signal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Signal(Debug, "channel_1", "channel 1 : debug message"))
			assert.NoError(t, facade.Signal(Fatal, "channel_2", "channel 2 : fatal message"))
			assert.NoError(t, facade.Signal(Fatal, "channel_1", "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		now := time.Now()
		fileName := fmt.Sprintf("/file-%s", now.Format("2006-01-02"))
		file, _ := disk.OpenFile(fileName, os.O_RDONLY, os.FileMode(0o644))
		data, _ := io.ReadAll(file)
		sdata := string(data)

		rx := `^`
		rx += `{\s*`
		rx += `"channel"\s*\:\s*"channel_1",\s*`
		rx += `"level"\s*\:\s*"FATAL",\s*`
		rx += `"message"\s*\:\s*"channel 1 : fatal message",\s*`
		rx += `"time"\s*\:\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}"\s*`
		rx += `}\s*`
		rx += `$`
		assert.Regexp(t, rx, sdata)
	})

	t.Run("should correctly handle the stream signal (any channel)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"*"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Signal(Debug, "channel_1", "channel 1 : debug message"))
			assert.NoError(t, facade.Signal(Fatal, "channel_2", "channel 2 : fatal message"))
			assert.NoError(t, facade.Signal(Fatal, "channel_1", "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		now := time.Now()
		fileName := fmt.Sprintf("/file-%s", now.Format("2006-01-02"))
		file, _ := disk.OpenFile(fileName, os.O_RDONLY, os.FileMode(0o644))
		data, _ := io.ReadAll(file)
		sdata := string(data)

		rx := `^`
		rx += `{\s*`
		rx += `"channel"\s*\:\s*"channel_2",\s*`
		rx += `"level"\s*\:\s*"FATAL",\s*`
		rx += `"message"\s*\:\s*"channel 2 : fatal message",\s*`
		rx += `"time"\s*\:\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}"\s*`
		rx += `}\s*`
		rx += `{\s*`
		rx += `"channel"\s*\:\s*"channel_1",\s*`
		rx += `"level"\s*\:\s*"FATAL",\s*`
		rx += `"message"\s*\:\s*"channel 1 : fatal message",\s*`
		rx += `"time"\s*\:\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}"\s*`
		rx += `}\s*`
		rx += `$`
		assert.Regexp(t, rx, sdata)
	})

	t.Run("should correctly handle the stream broadcast", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, NewProvider().Register(container))

		disk := afero.NewMemMapFs()

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)

		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Broadcast(Debug, "channel 1 : debug message"))
			assert.NoError(t, facade.Broadcast(Fatal, "channel 2 : fatal message"))
			assert.NoError(t, facade.Broadcast(Fatal, "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		now := time.Now()
		fileName := fmt.Sprintf("/file-%s", now.Format("2006-01-02"))
		file, _ := disk.OpenFile(fileName, os.O_RDONLY, os.FileMode(0o644))
		data, _ := io.ReadAll(file)
		sdata := string(data)

		rx := `^`
		rx += `{\s*`
		rx += `"level"\s*\:\s*"FATAL",\s*`
		rx += `"message"\s*\:\s*"channel 2 : fatal message",\s*`
		rx += `"time"\s*\:\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}"\s*`
		rx += `}\s*`
		rx += `{\s*`
		rx += `"level"\s*\:\s*"FATAL",\s*`
		rx += `"message"\s*\:\s*"channel 1 : fatal message",\s*`
		rx += `"time"\s*\:\s*"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}"\s*`
		rx += `}\s*`
		rx += `$`
		assert.Regexp(t, rx, sdata)
	})

	t.Run("should return the stream writer closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		provider := NewProvider()
		require.NoError(t, flamTime.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		expectedError := errors.New("expected error")
		file := NewFileMock(ctrl)
		file.EXPECT().Close().Return(expectedError).Times(1)

		now := time.Now()
		fileName := fmt.Sprintf("/file-%s", now.Format("2006-01-02"))
		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(file, nil).
			Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)

		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			provider.(flam.ClosableProvider).Close(container),
			expectedError)
	})

	t.Run("should correctly rotate the file on date change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		timeFacade := NewTimeFacadeMock(ctrl)
		gomock.InOrder(
			timeFacade.EXPECT().Now().Return(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeFacade.EXPECT().Now().Return(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
		)
		provider := NewProvider()
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		file1 := NewFileMock(ctrl)
		file1.EXPECT().Close().Return(nil).Times(1)
		file2 := NewFileMock(ctrl)
		file2.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)
		file2.EXPECT().Close().Return(nil).Times(1)

		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/file-2021-01-01", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(file1, nil).
			Times(1)
		disk.EXPECT().
			OpenFile("/file-2021-01-02", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(file2, nil).
			Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)

		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Broadcast(Fatal, "debug message"))
			assert.NoError(t, facade.Flush())
		}))

		assert.NoError(t, provider.(flam.ClosableProvider).Close(container))
	})

	t.Run("should return the rotating file opening error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(filesystem.PathDisks, flam.Bag{
			"mock": flam.Bag{
				"driver": "mock",
			}})
		_ = config.Defaults.Set(PathSerializers, flam.Bag{
			"json": flam.Bag{
				"driver": SerializerDriverJson,
			}})
		_ = config.Defaults.Set(PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     StreamDriverRotatingFile,
				"serializer": "json",
				"level":      "warning",
				"disk":       "mock",
				"path":       "/file-%s",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		timeFacade := NewTimeFacadeMock(ctrl)
		gomock.InOrder(
			timeFacade.EXPECT().Now().Return(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			timeFacade.EXPECT().Now().Return(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
		)
		provider := NewProvider()
		require.NoError(t, container.Provide(func() flamTime.Facade { return timeFacade }))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, provider.Register(container))

		file1 := NewFileMock(ctrl)

		expectedError := errors.New("expected error")
		disk := NewDiskMock(ctrl)
		disk.EXPECT().
			OpenFile("/file-2021-01-01", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(file1, nil).
			Times(1)
		disk.EXPECT().
			OpenFile("/file-2021-01-02", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644)).
			Return(nil, expectedError).
			Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock", "driver": "mock"}
		diskCreator := NewDiskCreatorMock(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)

		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, provider.(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NoError(t, facade.Broadcast(Fatal, "debug message"))
			assert.ErrorIs(t, facade.Flush(), expectedError)
		}))
	})
}
