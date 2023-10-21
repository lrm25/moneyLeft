package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// Log levels
const (
	LevelCrit  = 1
	LevelError = 2
	LevelWarn  = 3
	LevelInfo  = 4
	LevelDebug = 5
	LevelTrace = 6
)

func getLevelString(level int) string {
	switch level {
	case LevelCrit:
		return "CRIT"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	}
	panic(fmt.Sprintf("Invalid log level %d", level))
}

func getLevel(level string) int {
	switch strings.ToUpper(level) {
	case "CRIT":
		return LevelCrit
	case "ERROR":
		return LevelError
	case "WARN":
		return LevelWarn
	case "INFO":
		return LevelInfo
	case "DEBUG":
		return LevelDebug
	case "TRACE":
		return LevelTrace
	}
	panic(fmt.Sprintf("Invalid log level string %s", level))
}

var logger *Logger

// Logger struct, with log level
type Logger struct {
	level int
	stream io.Writer
}

// Init initializes logger with a certain log level
func Init(level int) *Logger {
	if logger == nil {
		logger = &Logger{
			level: level,
			stream: os.Stdout,
		}
	}
	return logger
}

// InitWithWriter initializes a logger with a specific writer
func InitWithWriter(level int, stream io.Writer) *Logger {
	if logger == nil {
		logger = &Logger{
			level: level,
			stream: stream,
		}
	}
	return logger
}

// Get returns the logger
func Get() *Logger {
	return logger
}

// SetLevel changes an already-initialized logger's level
func (l *Logger) SetLevel(level string) {
	l.level = getLevel(level)
}

// Log a message.  If the level integer is above the level the logger is set to, ignore the message.
func (l *Logger) Log(level int, message string) {
	if l == nil {
		l = Init(LevelInfo)
	}
	if level <= l.level {
		datetime := time.Now().Format("2006-01-02 15:04:05 -0700")
		_, file, line, _ := runtime.Caller(2)
		fmt.Fprintf(l.stream, "%s %s [%s:%d] %s\n", datetime, getLevelString(level), file, line, message)
	}
}

// Trace message
func (l *Logger) Trace(message string) {
	l.Log(LevelTrace, message)
}

// Debug message
func (l *Logger) Debug(message string) {
	l.Log(LevelDebug, message)
}

// Info message
func (l *Logger) Info(message string) {
	l.Log(LevelInfo, message)
}

// Warn message
func (l *Logger) Warn(message string) {
	l.Log(LevelWarn, message)
}

// Error message
func (l *Logger) Error(message string) {
	l.Log(LevelError, message)
}

// Crit message - log and terminate
func (l *Logger) Crit(message string) {
	l.Log(LevelCrit, message)
	os.Exit(1)
}