package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	LEVEL_CRIT  = 1
	LEVEL_ERROR = 2
	LEVEL_WARN  = 3
	LEVEL_INFO  = 4
	LEVEL_DEBUG = 5
	LEVEL_TRACE = 6
)

func getLevelString(level int) string {
	switch level {
	case LEVEL_CRIT:
		return "CRIT"
	case LEVEL_ERROR:
		return "ERROR"
	case LEVEL_WARN:
		return "WARN"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_TRACE:
		return "TRACE"
	}
	panic(fmt.Sprintf("Invalid log level %d", level))
}

func getLevel(level string) int {
	switch strings.ToUpper(level) {
	case "CRIT":
		return LEVEL_CRIT
	case "ERROR":
		return LEVEL_ERROR
	case "WARN":
		return LEVEL_WARN
	case "INFO":
		return LEVEL_INFO
	case "DEBUG":
		return LEVEL_DEBUG
	case "TRACE":
		return LEVEL_TRACE
	}
	panic(fmt.Sprintf("Invalid log level string %s", level))
}

var logger *Logger = nil

type Logger struct {
	level int
}

func Init(level int) *Logger {
	if logger == nil {
		logger = &Logger{
			level: level,
		}
	}
	return logger
}

func Get() *Logger {
	return logger
}

func (l *Logger) SetLevel(level string) {
	l.level = getLevel(level)
}

func (l *Logger) Log(level int, message string) {
	if l == nil {
		l = Init(LEVEL_INFO)
	}
	if level <= l.level {
		datetime := time.Now().Format("2006-01-02 15:04:05 -0700")
		_, file, line, _ := runtime.Caller(2)
		fmt.Printf("%s %s [%s:%d] %s\n", datetime, getLevelString(level), file, line, message)
	}
}

func (l *Logger) Trace(message string) {
	l.Log(LEVEL_TRACE, message)
}

func (l *Logger) Debug(message string) {
	l.Log(LEVEL_DEBUG, message)
}

func (l *Logger) Info(message string) {
	l.Log(LEVEL_INFO, message)
}

func (l *Logger) Warn(message string) {
	l.Log(LEVEL_WARN, message)
}

func (l *Logger) Error(message string) {
	l.Log(LEVEL_ERROR, message)
}

func (l *Logger) Crit(message string) {
	l.Log(LEVEL_CRIT, message)
	os.Exit(1)
}