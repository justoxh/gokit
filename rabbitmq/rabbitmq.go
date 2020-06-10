package rabbitmq


type (
	Config struct {
		Username string
		Password string
		Host string
		Port int
	}

	QueueConfig struct {
		Name       string `json:"name"`
		RoutingKey string `json:"routing_key"`
	}

	ExchangeConfig struct {
		Name   string        `json:"name"`
		Type   string        `json:"type"`
		Queues []QueueConfig `json:"queues"`
	}
)

