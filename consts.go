package log

const (
	providerId = "flam.log.provider"

	SerializerCreatorGroup   = "flam.log.serializers.creator"
	SerializerDriverString   = "flam.log.serializers.driver.string"
	SerializerDriverJson     = "flam.log.serializers.driver.json"
	StreamCreatorGroup       = "flam.log.streams.creator"
	StreamDriverConsole      = "flam.log.streams.driver.console"
	StreamDriverFile         = "flam.log.streams.driver.file"
	StreamDriverRotatingFile = "flam.log.streams.driver.rotating-file"

	PathDefaultLevel      = "flam.log.defaults.level"
	PathDefaultSerializer = "flam.log.defaults.serializer"
	PathDefaultDisk       = "flam.log.defaults.disk"
	PathBoot              = "flam.log.boot"
	PathFlusherFrequency  = "flam.log.flusher.frequency"
	PathSerializers       = "flam.log.serializers"
	PathStreams           = "flam.log.streams"
)
