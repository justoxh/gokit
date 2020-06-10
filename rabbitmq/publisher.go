package rabbitmq

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	con         *amqp.Connection
	channel     *amqp.Channel
	msgCh       chan *Msg
	reconnectCh chan struct{}
	wg          *sync.WaitGroup
}

type Msg struct {
	exchange   string
	routingKey string
	body       []byte
	header     map[string]interface{}
}

var DefaultPublisher *Publisher

func NewPublisher(amqpURI string) (*Publisher, error) {
	//amqp://guest:guest@localhost:5672/

	p := &Publisher{
		msgCh:       make(chan *Msg, 10000),
		reconnectCh: make(chan struct{}, 1),
	}

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}
	p.con = conn
	p.HandleReconnect()
	keepAlive(conn, amqpURI, &p.con, p.reconnectCh)

	channel, err := p.con.Channel()
	if err != nil {

		return nil, err
	}
	p.channel = channel
	p.Run()

	DefaultPublisher = p
	return DefaultPublisher, nil
}

func (p *Publisher) Close() {
	if p.con != nil {
		close(p.msgCh)
		p.con.Close()
		p.wg.Wait()
	}
}

func (c *Publisher) HandleReconnect() {
	go func() {
		for {
			<-c.reconnectCh
			channel, err := c.con.Channel()
			if err != nil {
				fmt.Printf("Run:%v\n", err)
				return
			}
			c.channel = channel
		}
	}()
}

func (p *Publisher) Run() {

	p.wg = &sync.WaitGroup{}
	p.wg.Add(1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
			if p.channel != nil {
				p.channel.Close()
			}
			p.wg.Done()
		}()

		for msg := range p.msgCh {
			pub := amqp.Publishing{
				Body:        msg.body,
				ContentType: "application/json",
				Headers:     msg.header,
			}
			for {
				fmt.Printf("Publish: exchange=%s, routingKey=%s, body=%s\n", msg.exchange, msg.routingKey, string(msg.body))
				err := p.channel.Publish(msg.exchange, msg.routingKey, false, false, pub)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				break
			}
		}
	}()
}

func (p *Publisher) Publish(exchange, routingKey string, body []byte, header map[string]interface{}) error {
	retry := 3 //错误重试
	timeWait := time.Second
	var err error
	for i := 0; i < retry; i++ {
		msg := amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
			Headers:     header,
		}
		fmt.Printf("Publish: exchange=%s, routingKey=%s, body=%s\n", exchange, routingKey, string(body))
		err = p.channel.Publish(exchange, routingKey, false, false, msg)

		if err != nil {
			time.Sleep(timeWait)
			continue
		}
		break
	}

	return err
}

func (p *Publisher) AsyncPublish(exchange, routingKey string, body []byte, header map[string]interface{}) error {
	msg := &Msg{exchange: exchange,
		routingKey: routingKey,
		body:       body,
		header:     header}

	p.msgCh <- msg
	return nil
}

func Publish(exchange, routingKey string, body []byte, header map[string]interface{}) error {
	return DefaultPublisher.Publish(exchange, routingKey, body, header)
}

func PublishJson(exchange, routingKey string, payload interface{}, header map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return DefaultPublisher.Publish(exchange, routingKey, body, header)
}

func AsyncPublish(exchange, routingKey string, body []byte, header map[string]interface{}) error {
	return DefaultPublisher.AsyncPublish(exchange, routingKey, body, header)
}

func AsyncPublishJson(exchange, routingKey string, payload interface{}, header map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return DefaultPublisher.AsyncPublish(exchange, routingKey, body, header)
}
