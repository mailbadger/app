package queue

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
)

const key = "producer"

type Producer interface {
	Publish(topic string, body []byte) error
	Stop()
}

func NewProducer(host, port string) (Producer, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	config := nsq.NewConfig()

	return nsq.NewProducer(addr, config)
}

// SetToContext sets the producer to the context
func SetToContext(c *gin.Context, producer Producer) {
	c.Set(key, producer)
}

// GetFromContext returns the Producer associated with the context
func GetFromContext(c context.Context) Producer {
	return c.Value(key).(Producer)
}

func Publish(c context.Context, topic string, body []byte) error {
	return GetFromContext(c).Publish(topic, body)
}
