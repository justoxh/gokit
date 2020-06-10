package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)



type Declarer struct {
	con *amqp.Connection
}

var DefaultDeclarer *Declarer

func NewDeclarer(amqpURI string) (*Declarer, error) {
	//amqp://guest:guest@localhost:5672/
	DefaultDeclarer = &Declarer{}

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}
	DefaultDeclarer.con = conn
	return DefaultDeclarer, nil
}

func (p *Declarer) Close() {
	if p.con != nil {
		p.con.Close()
	}
}

func (d *Declarer) Declare(exchanges []ExchangeConfig) error {
	channel, err := d.con.Channel()
	if err != nil {
		return fmt.Errorf("Declare: %s\n", err)
	}
	defer channel.Close()

	for _, ex := range exchanges {
		if ex.Name != "" && ex.Name != "default" {
			if err := d.ExchangeDeclare(channel, ex.Name, ex.Type); err != nil {
				return err
			}
		}
		for _, q := range ex.Queues {
			if err := d.QueueDeclare(channel, q.Name, true); err != nil {
				return err
			}
			if ex.Name != "" && ex.Name != "default" {
				if err := d.QueueBind(channel, q.Name, q.RoutingKey, ex.Name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *Declarer) ExchangeDeclare(channel *amqp.Channel, exchange, exchangeType string) error {
	err := channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *Declarer) QueueDeclare(channel *amqp.Channel, queueName string, durable bool) error {
	_, err := channel.QueueDeclare(queueName, durable, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *Declarer) QueueBind(channel *amqp.Channel, queueName, routingkey, exchange string) error {
	err := channel.QueueBind(queueName, routingkey, exchange, false, amqp.Table{})
	if err != nil {
		return err
	}
	return nil
}

func (d *Declarer) QueueUnbind(channel *amqp.Channel, queueName, routingkey, exchange string) error {
	err := channel.QueueUnbind(queueName, routingkey, exchange, amqp.Table{})
	if err != nil {
		return err
	}
	return nil
}
