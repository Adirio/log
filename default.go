package log

import (
	"fmt"
	"io"
)

var defaultLogger *Logger

func init() {
	defaultLogger = New(FLevel | FDate | FTime | FMicroseconds | FFile | FPath | FLine)
}

/************************\
|* Flag related methods *|
\************************/

func Set(flag flag) {
	defaultLogger.Set(flag)
}

func Unset(flag flag) {
	defaultLogger.Unset(flag)
}

/**************************\
|* Output related methods *|
\**************************/

func AddOutput (alias string, w io.Writer, level Level) error {
	return defaultLogger.AddOutput(alias, w, level)
}

func UpdateOutput(alias string, w io.Writer, level Level) error {
	return defaultLogger.UpdateOutput(alias, w, level)
}

func ReplaceOutput(alias string, w io.Writer, level Level) error {
	return defaultLogger.ReplaceOutput(alias, w, level)
}

func RemoveOutput(alias string) error {
	return defaultLogger.RemoveOutput(alias)
}

func ClearOutputs() {
	defaultLogger.ClearOutputs()
}

/*************************\
|* Level related methods *|
\*************************/

func GetLevel(alias string) (Level, error) {
	return defaultLogger.GetLevel(alias)
}

func SetLevel(alias string, level Level) error {
	return defaultLogger.SetLevel(alias, level)
}

/*************************\
|* Write related methods *|
\*************************/
// They can not use the logger exported API as the depth would be increased

func Trace(v ...interface{}) ([]int, []error) {
	return defaultLogger.trace(fmt.Sprint(v...))
}

func Tracef(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.trace(fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) ([]int, []error) {
	return defaultLogger.debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) ([]int, []error) {
	return defaultLogger.info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.info(fmt.Sprintf(format, v...))
}

func Notice(v ...interface{}) ([]int, []error) {
	return defaultLogger.notice(fmt.Sprint(v...))
}

func Noticef(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.notice(fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) ([]int, []error) {
	return defaultLogger.warning(fmt.Sprint(v...))
}

func Warningf(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.warning(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) ([]int, []error) {
	return defaultLogger.error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) ([]int, []error) {
	return defaultLogger.error(fmt.Sprintf(format, v...))
}

func Critical(v ...interface{}) {
	defaultLogger.critical(fmt.Sprint(v...))
}

func Criticalf(format string, v ...interface{}) {
	defaultLogger.critical(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	defaultLogger.fatal(fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.fatal(fmt.Sprintf(format, v...))
}
