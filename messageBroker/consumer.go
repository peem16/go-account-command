package messageBroker

import (
	"github.com/Shopify/sarama"
)

type EventHandler interface {
	Handle(topic string, eventBytes []byte)
}

type consumerHandler struct {
	eventHandler EventHandler
}

func NewConsumerHandler(eventHandler EventHandler) sarama.ConsumerGroupHandler {
	return consumerHandler{eventHandler: eventHandler}
}

func (obj consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		obj.eventHandler.Handle(msg.Topic, msg.Value)
		session.MarkMessage(msg, "")
	}

	return nil
}
