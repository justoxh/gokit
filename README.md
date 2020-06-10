## logrus log

### install
  
```
go get github.com/justoxh/gokit/logger
```
### usage
  
```
cfg := &logger.Options{
    Formatter:      "json",
    Write:          false,
    Path:           "./logs/",
    DisableConsole: false,
    WithCallerHook: true,
    MaxAge:         time.Duration(7*24) * time.Hour,
    RotationTime:   time.Duration(7) * time.Hour,
}
log := logger.GetLoggerWithOptions("default", cfg)
log.Info("Hello world")
```
  
  
## rabbitmq

### install 

```
go get github.com/justoxh/gokit/rabbitmq
```

### usage

```go
cfg := &rabbitmq.Config{
    Username: "guest",
    Password: "guest",
    Host:     "127.0.0.1",
    Port:     5672,
}
uri := rabbitmq.MakeURI(cfg)

declarer, _ := rabbitmq.NewDeclarer(uri)

publisher, _ := rabbitmq.NewPublisher(uri)

consumer, _ := rabbitmq.NewConsumer(uri, queue, "", 0, worker)

```



