package rabbitmq

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/streadway/amqp"
)

type Handler interface {
	HandelMessage(body []byte) error
}

type Consumer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	tag         string
	wg          *sync.WaitGroup
	queueName   string
	prefetch    int
	msgHandler  Handler
	reconnectCh chan struct{}
}

func MakeURI(cfg *Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

func NewConsumer(amqpURI, queueName, ctag string, prefetch int, msgHandler Handler) (*Consumer, error) {
	c := &Consumer{
		conn:        nil,
		channel:     nil,
		tag:         ctag,
		queueName:   queueName,
		prefetch:    prefetch,
		msgHandler:  msgHandler,
		reconnectCh: make(chan struct{}, 1),
	}

	var err error

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	c.HandleReconnect()
	keepAlive(conn, amqpURI, &c.conn, c.reconnectCh)
	err = c.Consume()

	return c, nil
}

func (c *Consumer) HandleReconnect() {
	go func() {
		for {
			<-c.reconnectCh
			c.Consume()
		}
	}()
}

func (c *Consumer) Consume() error {
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {

		return err
	}

	err = c.channel.Qos(c.prefetch, 0, false)
	if err != nil {

		return err
	}

	deliveries, err := c.channel.Consume(c.queueName, c.tag, false, false, false, false, nil)
	if err != nil {
		return err
	}
	num := 1
	if c.prefetch > 1 {
		num = c.prefetch
	}
	c.wg = &sync.WaitGroup{}
	for i := 0; i < num; i++ {
		c.wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
				}
				c.wg.Done()
			}()
			c.handle(deliveries)
		}()
	}
	return nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer func() {
		fmt.Println("AMQP shutdown OK")
	}()

	c.wg.Wait()
	return nil
}

func (c *Consumer) handle(deliveries <-chan amqp.Delivery) {
	fmt.Printf("deliveries count:%v\n", len(deliveries))
	for deliver := range deliveries {
		err := c.msgHandler.HandelMessage(deliver.Body)
		if err != nil {
			if err := deliver.Nack(false, true); err != nil {
				fmt.Printf("deliver.Nack: %s\n", err)
			}
		} else {
			if err := deliver.Ack(false); err != nil {
				fmt.Printf("deliver.Ack: %s\n", err)
			}
		}

	}
	fmt.Println("handle: deliveries channel closed")
}
