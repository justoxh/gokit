package logger

import (
	"testing"
	"time"
)

func TestGetLoggerWithOptions(t *testing.T) {
	options := &Options{
		Formatter:      "json",
		Write:          false,
		Path:           "./logs/",
		DisableConsole: false,
		WithCallerHook: true,
		MaxAge:         time.Duration(7*24) * time.Hour,
		RotationTime:   time.Duration(7) * time.Hour,
	}

	log := GetLoggerWithOptions("default", options)
	log.Info("Hello world")
}
