package logging

import (
	"fmt"
	"os"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level  Level
	output *os.File
}

func NewLogger(level Level) *Logger {
	return &Logger{level: level, output: os.Stdout}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.log("ERROR", msg, args...)
	}
}

func (l *Logger) log(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	fmt.Fprintf(l.output, "[%s] %s: %s\n", timestamp, level, formatted)
}
