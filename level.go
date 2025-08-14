package log

type Level int

const (
	None Level = iota
	Fatal
	Error
	Warning
	Notice
	Info
	Debug
)

var LevelName = map[Level]string{
	None:    "none",
	Fatal:   "fatal",
	Error:   "error",
	Warning: "warning",
	Notice:  "notice",
	Info:    "info",
	Debug:   "debug",
}

var LevelMap = map[string]Level{
	"none":    None,
	"fatal":   Fatal,
	"error":   Error,
	"warning": Warning,
	"notice":  Notice,
	"info":    Info,
	"debug":   Debug,
}

func LevelFrom(
	val any,
	def ...Level,
) Level {
	switch v := val.(type) {
	case Level:
		return v
	case int:
		if v >= int(None) && v <= int(Debug) {
			return Level(v)
		} else if len(def) != 0 {
			return def[0]
		}
	case string:
		if level, ok := LevelMap[v]; ok {
			return level
		} else if len(def) != 0 {
			return def[0]
		}
	default:
		if len(def) != 0 {
			return def[0]
		}
	}

	return None
}
