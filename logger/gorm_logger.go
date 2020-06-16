package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type GormLogger struct {
	Log *logrus.Logger
}

var gormLogger *GormLogger

func (logger *GormLogger) Print(v ...interface{}){
	switch v[0] {
	case "sql":
		logger.Log.WithFields(logrus.Fields{
			"type":          "sql",
			"rows_returned": v[5],
			"src":           v[1],
			"values":        v[4],
			"duration":      fmt.Sprintf("%.2fms",float64(v[2].(time.Duration).Nanoseconds()/1e4)/100.0),
		},
		).Info(v[3])
	case "log":
		logger.Log.WithFields(logrus.Fields{ "type": "log"}).Print(v[2])
	}
}



func GetGormLoggerWithOptions(options *Options) *GormLogger {
	var logger GormLogger
	logger.Log = GetLoggerWithOptions("gorm",options)
	gormLogger =&logger
	return gormLogger
}