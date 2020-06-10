package rabbitmq

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func makeURI(cfg *Config) string {
	return MakeURI(cfg)
}

func config() *Config {
	return &Config{
		Username: "guest",
		Password: "guest",
		Host:     "127.0.0.1",
		Port:     5672,
	}
}

func exchanges() []ExchangeConfig {
	return []ExchangeConfig{
		{
			Name: "order",
			Type: "direct",
			Queues: []QueueConfig{
				{
					Name:       "create",
					RoutingKey: "order_create",
				},
				{
					Name:       "update",
					RoutingKey: "order_update",
				},
			},
		},
	}
}

func TestNewDeclarer(t *testing.T) {
	cfg := config()
	uri := makeURI(cfg)
	declarer, err := NewDeclarer(uri)
	if err != nil {
		t.Logf("NewDeclarer:%v", err)
	}
	defer declarer.Close()
	if err := declarer.Declare(exchanges()); err != nil {
		t.Logf("Declare:%v", err)
	}
	t.Log("success")
}

const (
	ORDER_CREATE = iota + 1
	ORDER_UPDATE
)

type Payload struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type OrderCreate struct {
	OrderId int64 `json:"order_id"`
	UserId  int64 `json:"user_id"`
}

type Worker struct {
	queue   string
}

func NewWorker(queue string) *Worker {
	return &Worker{
		queue:   queue,
	}
}

func (w *Worker) HandelMessage(body []byte) error {
	payload := Payload{}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}
	log.Println(payload)
	err = w.process(payload)
	if err != nil {
		return err
	}
	return nil
}
func (w *Worker) process(payload Payload) error {
	switch payload.Type {
	case ORDER_CREATE:

	}
	return nil
}

func TestNewPublisher(t *testing.T) {
	cfg := config()
	uri := makeURI(cfg)
	publisher, err := NewPublisher(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()
	message := OrderCreate{
		OrderId: 1,
		UserId:  1,
	}
	data, err := json.Marshal(&message)
	if err != nil {
		t.Fatal(err)
	}
	payload := Payload{
		Type: ORDER_CREATE,
		Data: string(data),
	}
	body, err := json.Marshal(&payload)
	if err != nil {
		t.Fatal(err)
	}
	header := amqp.Table{}
	if err := publisher.Publish("order", "order_create", body, header);err != nil {
		t.Fatal(err)
	}

}

func TestNewConsumer(t *testing.T) {
	cfg := config()
	uri := makeURI(cfg)
	exchanges := exchanges()
	var consumers []*Consumer
	for _, exchange := range exchanges {
		for _, queue := range exchange.Queues {
			worker := NewWorker(queue.Name)
			consumer, err := NewConsumer(uri, queue.Name,
				"", 0, worker)
			if err != nil {
				t.Fatal(err)
			}
			consumers = append(consumers,consumer)
		}
	}
	time.Sleep(time.Second * 2)
	for _,v := range consumers {
		_ = v.Shutdown()
	}
}


