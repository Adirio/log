package log

import "errors"

type Level int

const (
	LevelAll Level = iota

	// Normal behaviour
	LevelTrace
	LevelDebug
	LevelInfo
	LevelNotice

	// Faulty behaviour
	LevelWarning
	LevelError
	LevelCritical
	LevelFatal

	LevelNone
)

func (l Level) String() string {
	switch l {
	case LevelAll:
		return "all"
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelNotice:
		return "notice"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	case LevelFatal:
		return "fatal"
	case LevelNone:
		return "none"
	default:
		return "unknown"
	}
}

var (
	maxLevelLength int
	levelError = errors.New("invalid level value")
)

func init() {
	for l := LevelAll; l <= LevelNone; l++ {
		if maxLevelLength < len(l.String()) {
			maxLevelLength = len(l.String())
		}
	}
}
