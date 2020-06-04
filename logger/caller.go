package logger

import (
	"github.com/sirupsen/logrus"
)

type CallerHook struct {
	module string
}

// Fire adds a caller field in logger instance
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data["module"] = hook.module
	return nil
}

// Levels returns supported levels
func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
