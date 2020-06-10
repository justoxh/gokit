package rabbitmq

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/streadway/amqp"
)

//amqp断开自动重连
func keepAlive(conn *amqp.Connection, amqpURI string, ppConn **amqp.Connection, reconnectCh chan struct{}) {
	go func() {
		defer func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("keepAlive: %v\n", err)
					debug.PrintStack()
				}
			}()
		}()

		err := <-conn.NotifyClose(make(chan *amqp.Error))

		if err != nil {
			fmt.Printf("keepAlive: %s\n", err.Error())
			for {
				newcon, err := amqp.Dial(amqpURI)
				if err == nil {
					*ppConn = newcon
					if reconnectCh != nil {
						reconnectCh <- struct{}{}
					}
					keepAlive(newcon, amqpURI, ppConn, reconnectCh)
					break
				}
				time.Sleep(time.Second)
			}
		} else {
			fmt.Println("keepAlive: rabbitmq connection closing")
		}
	}()
}
