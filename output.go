package log

import (
	"io"
	"sync"
)

type output struct {
	mu    sync.RWMutex
	w     io.Writer
	level Level
}

func (out output) getLevel() Level {
	out.mu.RLock()
	defer out.mu.RUnlock()

	return out.level
}

func (out *output) setLevel(level Level) error{
	if level < LevelAll || level > LevelNone {
		return levelError
	}

	out.mu.Lock()
	defer out.mu.Unlock()

	out.level = level
	return nil
}

func (out output) write(bytes []byte, level Level) (int, error) {
	out.mu.Lock()
	defer out.mu.Unlock()

	if level < out.level {
		return 0, nil
	}

	return out.w.Write(bytes)
}
