## logrus log

### install
  
```
go get github.com/justoxh/gokit/logger
```
### usege
  
```
options := &logger.Options{
    Formatter:      "json",
    Write:          false,
    Path:           "./logs/",
    DisableConsole: false,
    WithCallerHook: true,
    MaxAge:         time.Duration(7*24) * time.Hour,
    RotationTime:   time.Duration(7) * time.Hour,
}
log := logger.GetLoggerWithOptions("default", options)
log.Info("Hello world")
```
  


