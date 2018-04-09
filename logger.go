package log

import (
	"github.com/pkg/errors"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"
	"fmt"
	"os"
)

var (
	aliasUsed = errors.New("alias already in use")
	aliasNotUsed = errors.New("alias not in use")
)

type Logger struct {
	mu      sync.RWMutex
	flags   flag
	outputs map[string]*output
	buf     []byte
}

func New(flags flag) *Logger {
	return &Logger{
		flags  : flags,
		outputs: make(map[string]*output),
		buf    : make([]byte, 80),
	}
}

/************************\
|* Flag related methods *|
\************************/

func (l *Logger) Set(flag flag) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.flags |= flag
}

func (l *Logger) Unset(flag flag) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.flags &= ^flag
}

/**************************\
|* Output related methods *|
\**************************/

func (l Logger) AddOutput(alias string, w io.Writer, level Level) error {
	if level < LevelAll || level > LevelNone {
		return levelError
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, found := l.outputs[alias]; found {
		return aliasUsed
	}

	l.outputs[alias] = &output{w: w, level: level}
	return nil
}

func (l Logger) UpdateOutput(alias string, w io.Writer, level Level) error {
	if level < LevelAll || level > LevelNone {
		return levelError
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.outputs[alias] = &output{w: w, level: level}
	return nil
}

func (l Logger) ReplaceOutput(alias string, w io.Writer, level Level) error {
	if level < LevelAll || level > LevelNone {
		return levelError
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, found := l.outputs[alias]; !found {
		return aliasNotUsed
	}

	l.outputs[alias] = &output{w: w, level: level}
	return nil
}

func (l Logger) RemoveOutput(alias string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, found := l.outputs[alias]; !found {
		return aliasNotUsed
	}

	delete(l.outputs, alias)
	return nil
}

func (l *Logger) ClearOutputs() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.outputs = make(map[string]*output)
}

/*************************\
|* Level related methods *|
\*************************/

func (l Logger) GetLevel(alias string) (Level, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, found := l.outputs[alias]; !found {
		return LevelAll, aliasNotUsed
	}

	return l.outputs[alias].getLevel(), nil
}

func (l Logger) SetLevel(alias string, level Level) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, found := l.outputs[alias]; !found {
		return aliasNotUsed
	}

	l.outputs[alias].setLevel(level)
	return nil
}

/*************************\
|* Write related methods *|
\*************************/

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l Logger) log(s string, level Level) ([]int, []error) {
	t := time.Now()
	var file string
	var line int

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flags & FFile != 0 {
		// Release lock while getting caller info as it is expensive
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	l.buf = l.buf[:0]
	if l.flags & FLevel != 0 {
		l.buf = append(l.buf, strings.Title(level.String())...)
		l.buf = append(l.buf, ' ')
	}
	if l.flags & (FDate|FTime) != 0 {
		if l.flags & FUTC != 0 {
			t = t.UTC()
		}
		l.buf = append(l.buf, '(')
		if l.flags & FDate != 0 {
			year, month, day := t.Date()
			itoa(&l.buf, year, 4)
			l.buf = append(l.buf, '/')
			itoa(&l.buf, int(month), 2)
			l.buf = append(l.buf, '/')
			itoa(&l.buf, day, 2)
			if l.flags & FTime != 0 {
				l.buf = append(l.buf, ' ')
			}
		}
		if l.flags & FTime != 0 {
			hour, min, sec := t.Clock()
			itoa(&l.buf, hour, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, min, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, sec, 2)
			if l.flags & FMicroseconds != 0 {
				l.buf = append(l.buf, '.')
				itoa(&l.buf, t.Nanosecond()/1e3, 6)
			}
		}
		l.buf = append(l.buf, ')', ' ')
	}
	if l.flags & FFile != 0 {
		if l.flags & FPath == 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		l.buf = append(l.buf, file...)
		if l.flags & FLine != 0 {
			l.buf = append(l.buf, ':')
			itoa(&l.buf, line, -1)
		}
		l.buf = append(l.buf, ':', ' ')
	}
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	bytes := make([]int, 0, len(l.outputs))
	errs := make([]error, 0, len(l.outputs))
	for _, out := range(l.outputs) {
		n, err := out.write(l.buf, level)
		bytes = append(bytes, n)
		errs = append(errs, err)
	}
	return bytes, errs
}

func (l Logger) Trace(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelTrace)
}

func (l Logger) Tracef(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelTrace)
}

func (l Logger) Debug(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelDebug)
}

func (l Logger) Debugf(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelDebug)
}

func (l Logger) Info(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelInfo)
}

func (l Logger) Infof(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelInfo)
}

func (l Logger) Notice(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelNotice)
}

func (l Logger) Noticef(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelNotice)
}

func (l Logger) Warning(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelWarning)
}

func (l Logger) Warningf(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelWarning)
}

func (l Logger) Error(v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprint(v...), LevelError)
}

func (l Logger) Errorf(format string, v ...interface{}) ([]int, []error) {
	return l.log(fmt.Sprintf(format, v...), LevelError)
}

func (l Logger) Critical(v ...interface{}) {
	s := fmt.Sprint(v...)
	bytes, errs := l.log(s, LevelCritical)
	panic(CriticalError{s, bytes, errs})
}

func (l Logger) Criticalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	bytes, errs :=  l.log(s, LevelCritical)
	panic(CriticalError{s, bytes, errs})
}

func (l Logger) Fatal(v ...interface{}) {
	l.log(fmt.Sprint(v...), LevelFatal)
	os.Exit(1)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.log(fmt.Sprintf(format, v...), LevelFatal)
	os.Exit(1)
}