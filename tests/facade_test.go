package tests

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	config "github.com/happyhippyhippo/flam-config"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	log "github.com/happyhippyhippo/flam-log"
	mocks "github.com/happyhippyhippo/flam-log/tests/mocks"
	time "github.com/happyhippyhippo/flam-time"
)

func Test_facade_Signal(t *testing.T) {
	t.Run("should not send message to stream if not flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(0)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Signal(log.Info, "channel", "message"))
		}))
	})

	t.Run("should send message to stream if flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Signal(log.Info, "channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send message with context to stream if flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Signal(log.Info, "channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_Broadcast(t *testing.T) {
	t.Run("should not send message to stream if not flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message", flam.Bag{}).
			Return(nil).
			Times(0)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Broadcast(log.Info, "message"))
		}))
	})

	t.Run("should send message to stream if flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Broadcast(log.Info, "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send message with context to stream if flushed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.Broadcast(log.Info, "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_FatalSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Fatal, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.FatalSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Fatal, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.FatalSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_FatalBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Fatal, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.FatalBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Fatal, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.FatalBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_ErrorSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Error, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.ErrorSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Error, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.ErrorSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_ErrorBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Error, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.ErrorBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Error, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.ErrorBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_WarningSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Warning, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.WarningSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Warning, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.WarningSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_WarningBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Warning, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.WarningBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Warning, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.WarningBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_NoticeSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Notice, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Notice, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_NoticeBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Notice, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Notice, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_InfoSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.InfoSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Info, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.InfoSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_InfoBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.InfoBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.InfoBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_DebugSignal(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Debug, "channel", "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.DebugSignal("channel", "message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Debug, "channel", "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.DebugSignal("channel", "message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_DebugBroadcast(t *testing.T) {
	t.Run("should send appropriate leveled message to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Debug, "message", flam.Bag{}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.DebugBroadcast("message"))
			assert.NoError(t, facade.Flush())
		}))
	})

	t.Run("should send appropriate leveled message with context to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Debug, "message", flam.Bag{"key": "value"}).
			Return(nil).
			Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.DebugBroadcast("message", flam.Bag{"key": "value"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_Flush(t *testing.T) {
	t.Run("should return the signal message forward action to stream error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("expected error")
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Signal(gomock.Any(), log.Notice, "channel", "message3", flam.Bag{"key3": "value3"}).
			Return(expectedErr)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeSignal("channel", "message3", flam.Bag{"key3": "value3"}))
			assert.ErrorIs(t, facade.Flush(), expectedErr)
		}))
	})

	t.Run("should return the braodcast message forward action to stream error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("expected error")
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Notice, "message3", flam.Bag{"key3": "value3"}).
			Return(expectedErr)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.NoticeBroadcast("message3", flam.Bag{"key3": "value3"}))
			assert.ErrorIs(t, facade.Flush(), expectedErr)
		}))
	})

	t.Run("should send stored messages to stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Debug, "message1", flam.Bag{"key1": "value1"}).
			Return(nil)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Info, "message2", flam.Bag{"key2": "value2"}).
			Return(nil)
		stream.EXPECT().
			Broadcast(gomock.Any(), log.Notice, "message3", flam.Bag{"key3": "value3"}).
			Return(nil)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.DebugBroadcast("message1", flam.Bag{"key1": "value1"}))
			assert.NoError(t, facade.InfoBroadcast("message2", flam.Bag{"key2": "value2"}))
			assert.NoError(t, facade.NoticeBroadcast("message3", flam.Bag{"key3": "value3"}))
			assert.NoError(t, facade.Flush())
		}))
	})
}

func Test_facade_HasSerializer(t *testing.T) {
	t.Run("should return false if serializer does not exist", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.False(t, facade.HasSerializer("serializer"))
		}))
	})

	t.Run("should return true on known serializer in config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.True(t, facade.HasSerializer("serializer"))
		}))
	})

	t.Run("should return true on an added serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		serializer := mocks.NewSerializer(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddSerializer("serializer", serializer))

			assert.True(t, facade.HasSerializer("serializer"))
		}))
	})
}

func Test_facade_ListSerializers(t *testing.T) {
	t.Run("should return an empty list if no serializer has been registered", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.Empty(t, facade.ListSerializers())
		}))
	})

	t.Run("should return a list of registered serializers (added)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddSerializer("serializer1", mocks.NewSerializer(ctrl)))
			require.NoError(t, facade.AddSerializer("serializer2", mocks.NewSerializer(ctrl)))

			assert.ElementsMatch(t, []string{"serializer1", "serializer2"}, facade.ListSerializers())
		}))
	})

	t.Run("should return a list of registered serializers (in config)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer1": flam.Bag{}, "serializer2": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.ElementsMatch(t, []string{"serializer1", "serializer2"}, facade.ListSerializers())
		}))
	})

	t.Run("should return a list of registered serializers (in config and added)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer1": flam.Bag{}, "serializer3": flam.Bag{}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(2)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddSerializer("serializer2", mocks.NewSerializer(ctrl)))

			assert.ElementsMatch(
				t,
				[]string{"serializer1", "serializer2", "serializer3"},
				facade.ListSerializers())
		}))
	})
}

func Test_facade_GetSerializer(t *testing.T) {
	t.Run("should return an error if serializer does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetSerializer("serializer")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrUnknownResource)
		}))
	})

	t.Run("should return error when unable to generate the serializer instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetSerializer("serializer")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrInvalidResourceConfig)
		}))
	})

	t.Run("should return error originated when generating the serializer instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		expectedErr := errors.New("expected error")
		serializerCreatorConfig := flam.Bag{"id": "serializer", "driver": "mock"}
		serializerCreator := mocks.NewSerializerCreator(ctrl)
		serializerCreator.EXPECT().Accept(serializerCreatorConfig).Return(true).Times(1)
		serializerCreator.EXPECT().Create(serializerCreatorConfig).Return(nil, expectedErr).Times(1)
		require.NoError(t, container.Provide(func() log.SerializerCreator {
			return serializerCreator
		}, dig.Group(log.SerializerCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetSerializer("serializer")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, expectedErr)
		}))
	})

	t.Run("should return 'string' serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{"driver": log.SerializerDriverString}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetSerializer("serializer")
			assert.NotNil(t, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return 'json' serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{"driver": log.SerializerDriverJson}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetSerializer("serializer")
			assert.NotNil(t, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_facade_AddSerializer(t *testing.T) {
	t.Run("should return error on nil serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.ErrorIs(t, facade.AddSerializer("serializer", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should store a serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		serializer := mocks.NewSerializer(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddSerializer("serializer", serializer))

			got, e := facade.GetSerializer("serializer")
			assert.Same(t, serializer, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return error when trying to store an already existing serializer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		cfg := flam.Bag{"serializer": flam.Bag{"driver": "mock"}}
		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.EXPECT().Get(log.PathSerializers).Return(cfg).Times(1)
		require.NoError(t, container.Decorate(func(flam.FactoryConfig) flam.FactoryConfig {
			return factoryConfig
		}))

		serializer := mocks.NewSerializer(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.ErrorIs(
				t,
				facade.AddSerializer("serializer", serializer),
				flam.ErrDuplicateResource)
		}))
	})
}

func Test_Facade_HasStream(t *testing.T) {
	t.Run("should return false on unknown stream", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.False(t, facade.HasStream("unknown"))
		}))
	})

	t.Run("should return false if the stream is on config but not loaded", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade config.Facade) {
			assert.NoError(t, facade.Set(log.PathStreams, flam.Bag{
				"stream": flam.Bag{
					"driver": "mock",
				}}))
		}))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.False(t, facade.HasStream("stream"))
		}))
	})

	t.Run("should return true if the stream was added directly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_ = config.Defaults.Set(log.PathStreams, flam.Bag{
			"stream": flam.Bag{
				"driver": "mock",
			}})
		defer func() { config.Defaults = flam.Bag{} }()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.True(t, facade.HasStream("stream"))
		}))
	})
}

func Test_Facade_ListStreams(t *testing.T) {
	t.Run("should return an empty list if no streams were added", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.Empty(t, facade.ListStreams())
		}))
	})

	t.Run("should return a ordered list of streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream1 := mocks.NewStream(ctrl)
		stream2 := mocks.NewStream(ctrl)
		stream3 := mocks.NewStream(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("zulu", stream3))
			require.NoError(t, facade.AddStream("alpha", stream1))
			require.NoError(t, facade.AddStream("charlie", stream2))

			assert.ElementsMatch(
				t,
				[]string{"alpha", "charlie", "zulu"},
				facade.ListStreams())
		}))
	})
}

func Test_Facade_GetStream(t *testing.T) {
	t.Run("should return error on unknown stream", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			got, e := facade.GetStream("unknown")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, log.ErrStreamNotFound)
		}))
	})

	t.Run("should return requested stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			got, e := facade.GetStream("stream")
			assert.Same(t, stream, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_Facade_AddStream(t *testing.T) {
	t.Run("should return error on nil stream", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.ErrorIs(t, facade.AddStream("stream", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should return error on duplicate stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))
			assert.ErrorIs(t, facade.AddStream("stream", stream), log.ErrDuplicateStream)
		}))
	})

	t.Run("should add stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			got, e := facade.GetStream("stream")
			assert.Same(t, stream, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_Facade_RemoveStream(t *testing.T) {
	t.Run("should return error on invalid stream", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			assert.ErrorIs(t, facade.RemoveStream("stream"), log.ErrStreamNotFound)
		}))
	})

	t.Run("should return the error if the stream returns any on closing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		expectedErr := errors.New("close error")
		stream := mocks.NewStream(ctrl)
		stream.EXPECT().Close().Return(expectedErr).Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.ErrorIs(t, facade.RemoveStream("stream"), expectedErr)
		}))
	})

	t.Run("should remove stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream := mocks.NewStream(ctrl)
		stream.EXPECT().Close().Return(nil).Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.NoError(t, facade.RemoveStream("stream"))

			assert.False(t, facade.HasStream("stream"))
		}))
	})
}

func Test_Facade_RemoveAllStreams(t *testing.T) {
	t.Run("should return the error if the stream returns any on closing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		expectedErr := errors.New("close error")
		stream := mocks.NewStream(ctrl)
		stream.EXPECT().Close().Return(expectedErr)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream", stream))

			assert.ErrorIs(t, facade.RemoveAllStreams(), expectedErr)
		}))
	})

	t.Run("should remove all streams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, time.NewProvider().Register(container))
		require.NoError(t, filesystem.NewProvider().Register(container))
		require.NoError(t, config.NewProvider().Register(container))
		require.NoError(t, log.NewProvider().Register(container))

		stream1 := mocks.NewStream(ctrl)
		stream1.EXPECT().Close().Return(nil).Times(1)

		stream2 := mocks.NewStream(ctrl)
		stream2.EXPECT().Close().Return(nil).Times(1)

		assert.NoError(t, container.Invoke(func(facade log.Facade) {
			require.NoError(t, facade.AddStream("stream1", stream1))
			require.NoError(t, facade.AddStream("stream2", stream2))

			assert.NoError(t, facade.RemoveAllStreams())

			assert.Empty(t, facade.ListStreams())
		}))
	})
}
