package tests

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	log "github.com/happyhippyhippo/flam-log"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_consoleStream(t *testing.T) {
	t.Run("should return serialization creation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "invalid",
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))

		assert.ErrorIs(
			t,
			log.NewProvider().(flam.BootableProvider).Boot(container),
			flam.ErrUnknownResource)
	})

	t.Run("should correctly handle the stream level", func(t *testing.T) {
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
				"level":      "debug",
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetStream("my_stream")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, log.Debug, got.GetLevel())

			assert.NoError(t, got.SetLevel(log.Info))
			assert.Equal(t, log.Info, got.GetLevel())
		}))
	})

	t.Run("should correctly handle the channel list", func(t *testing.T) {
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
				"level":      "debug",
				"channels":   []any{"channel_2", "channel_1"},
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetStream("my_stream")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.True(t, got.HasChannel("channel_1"))
			assert.True(t, got.HasChannel("channel_2"))
			assert.False(t, got.HasChannel("channel_3"))
			assert.ElementsMatch(t, got.ListChannels(), []string{"channel_1", "channel_2"})

			require.NoError(t, got.AddChannel("channel_3"))
			assert.True(t, got.HasChannel("channel_1"))
			assert.True(t, got.HasChannel("channel_2"))
			assert.True(t, got.HasChannel("channel_3"))
			assert.ElementsMatch(t, got.ListChannels(), []string{"channel_1", "channel_2", "channel_3"})

			require.NoError(t, got.RemoveChannel("channel_2"))
			assert.True(t, got.HasChannel("channel_1"))
			assert.False(t, got.HasChannel("channel_2"))
			assert.True(t, got.HasChannel("channel_3"))
			assert.ElementsMatch(t, got.ListChannels(), []string{"channel_1", "channel_3"})

			require.NoError(t, got.RemoveAllChannels())
			assert.False(t, got.HasChannel("channel_1"))
			assert.False(t, got.HasChannel("channel_2"))
			assert.False(t, got.HasChannel("channel_3"))
			assert.ElementsMatch(t, got.ListChannels(), []string{})
		}))
	})

	t.Run("should correctly handle the stream signal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			outC <- buf.String()
		}()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
				"level":      "warning",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.Signal(log.Debug, "channel_1", "channel 1 : debug message"))
			assert.NoError(t, facade.Signal(log.Fatal, "channel_2", "channel 2 : fatal message"))
			assert.NoError(t, facade.Signal(log.Fatal, "channel_1", "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		_ = w.Close()
		os.Stdout = old // restoring the real stdout
		out := <-outC

		rx := `^`
		rx += `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}\s*`
		rx += `\[FATAL\]\s*`
		rx += `channel 1 : fatal message\s*`
		rx += `$`
		assert.Regexp(t, rx, out)
	})

	t.Run("should correctly handle the stream signal (any channel)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			outC <- buf.String()
		}()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
				"level":      "warning",
				"channels":   []any{"*"},
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.Signal(log.Debug, "channel_1", "channel 1 : debug message"))
			assert.NoError(t, facade.Signal(log.Fatal, "channel_2", "channel 2 : fatal message"))
			assert.NoError(t, facade.Signal(log.Fatal, "channel_1", "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		_ = w.Close()
		os.Stdout = old // restoring the real stdout
		out := <-outC

		rx := `^`
		rx += `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}\s*`
		rx += `\[FATAL\]\s*`
		rx += `channel 2 : fatal message\s*`
		rx += `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}\s*`
		rx += `\[FATAL\]\s*`
		rx += `channel 1 : fatal message\s*`
		rx += `$`
		assert.Regexp(t, rx, out)
	})

	t.Run("should correctly handle the broadcast signal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		config.Defaults = flam.Bag{}
		_ = config.Defaults.Set(log.PathSerializers, flam.Bag{
			"string": flam.Bag{
				"driver": log.SerializerDriverString,
			}})
		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"my_stream": flam.Bag{
				"driver":     log.StreamDriverConsole,
				"serializer": "string",
				"level":      "warning",
				"channels":   []any{"channel_1"},
			}})
		_ = config.Defaults.Set(log.PathBoot, true)
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		require.NoError(t, config.NewProvider().(flam.BootableProvider).Boot(container))
		require.NoError(t, log.NewProvider().(flam.BootableProvider).Boot(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.NoError(t, facade.Broadcast(log.Debug, "channel 1 : debug message"))
			assert.NoError(t, facade.Broadcast(log.Fatal, "channel 2 : fatal message"))
			assert.NoError(t, facade.Broadcast(log.Fatal, "channel 1 : fatal message"))
			assert.NoError(t, facade.Flush())
		}))

		w.Close()
		os.Stdout = old // restoring the real stdout
		out := <-outC

		rx := `^`
		rx += `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}\s*`
		rx += `\[FATAL\]\s*`
		rx += `channel 2 : fatal message\s*`
		rx += `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}\+\d{4}\s*`
		rx += `\[FATAL\]\s*`
		rx += `channel 1 : fatal message\s*`
		rx += `$`
		assert.Regexp(t, rx, out)
	})
}
