package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel       = "debug"
	defaultLogFileName = "daily"
)

var (
	logMap      map[string]*logrus.Logger
	getLogMutex sync.Mutex
)

var logLevelMap = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}

func GetLogLevel(logLevelName string) logrus.Level {
	if v, ok := logLevelMap[logLevelName]; ok {
		return v
	}
	return logrus.DebugLevel
}

func defaultOptions() *Options {
	return &Options{
		Level:          defaultLevel,
		WithCallerHook: true,
		Formatter:      "text",
		DisableConsole: false,
		Write:          false,
		Path:           os.TempDir(),
		FileName:       defaultLogFileName,
		MaxAge:         time.Duration(24*7) * time.Hour,
		RotationTime:   time.Duration(24) * time.Hour,
		Debug:          false,
	}
}

type Options struct {
	Level          string
	WithCallerHook bool
	Formatter      string
	DisableConsole bool
	Write          bool
	Path           string
	FileName       string
	MaxAge         time.Duration
	RotationCount  uint
	RotationTime   time.Duration
	Debug          bool
}

func GetLoggerWithOptions(logName string, options *Options) *logrus.Logger {
	getLogMutex.Lock()
	defer getLogMutex.Unlock()

	if options == nil {
		options = defaultOptions()
	}

	if logMap == nil {
		logMap = make(map[string]*logrus.Logger)
	}
	curLog, ok := logMap[logName]

	if ok {
		return curLog
	}

	log := logrus.New()

	level := options.Level
	if level == "" {
		level = defaultLevel
	}
	logLevel := GetLogLevel(level)
	logDir := options.Path
	if logDir == "" {
		logDir = os.TempDir()
	}

	logFileName := options.FileName
	if logFileName == "" {
		logFileName = defaultLogFileName
	}

	printLog := !options.DisableConsole
	maxAge := options.MaxAge
	rotationCount := options.RotationCount
	rotationTime := options.RotationTime
	withCallerHook := options.WithCallerHook

	log.SetLevel(logLevel)

	if options.Write {
		storeLogDir := logDir

		err := os.MkdirAll(storeLogDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("creating log file failed: %s", err.Error()))
		}
		path := filepath.Join(storeLogDir, logFileName)
		writer, err := rotatelogs.New(
			path+"-%Y-%m-%d.log",
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour),
			rotatelogs.WithRotationCount(rotationCount),
			rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),
		)
		if err != nil {
			panic(fmt.Sprintf("rotatelogs log failed: %s", err.Error()))
		}

		var formatter logrus.Formatter

		if options.Formatter == "json" {
			formatter = &logrus.JSONFormatter{}
		} else {
			formatter = &logrus.TextFormatter{}
		}
		log.SetFormatter(formatter)
		if options.Debug {
			log.AddHook(lfshook.NewHook(
				lfshook.WriterMap{
					logrus.DebugLevel: writer,
					logrus.InfoLevel:  writer,
					logrus.WarnLevel:  writer,
					logrus.ErrorLevel: writer,
					logrus.FatalLevel: writer,
				},
				formatter,
			))

			defaultLogFilePrefix := logFileName + "."
			pathMap := lfshook.PathMap{
				logrus.DebugLevel: fmt.Sprintf("%s/%sdebug", storeLogDir, defaultLogFilePrefix),
				logrus.InfoLevel:  fmt.Sprintf("%s/%sinfo", storeLogDir, defaultLogFilePrefix),
				logrus.WarnLevel:  fmt.Sprintf("%s/%swarn", storeLogDir, defaultLogFilePrefix),
				logrus.ErrorLevel: fmt.Sprintf("%s/%serror", storeLogDir, defaultLogFilePrefix),
				logrus.FatalLevel: fmt.Sprintf("%s/%sfatal", storeLogDir, defaultLogFilePrefix),
			}
			log.AddHook(lfshook.NewHook(
				pathMap,
				formatter,
			))
		} else {
			log.Out = writer
		}

	} else {
		if printLog {
			log.Out = os.Stdout
			var formatter logrus.Formatter
			if options.Formatter == "json" {
				formatter = &logrus.JSONFormatter{}
				(formatter).(*logrus.JSONFormatter).TimestampFormat = "2006-01-02 15:04:05"
				(formatter).(*logrus.JSONFormatter).DisableTimestamp = false
			} else {
				formatter = &logrus.TextFormatter{}
				(formatter).(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05"
				(formatter).(*logrus.TextFormatter).FullTimestamp = true
				(formatter).(*logrus.TextFormatter).DisableTimestamp = false
			}
			log.Formatter = formatter
		}
	}

	if withCallerHook {
		log.AddHook(&CallerHook{module: logName})
	}
	curLog = log
	logMap[logName] = curLog
	fmt.Printf("register logger %v, store in %v, current loggers: %v\n", logName, logDir, logMap)
	return curLog
}
