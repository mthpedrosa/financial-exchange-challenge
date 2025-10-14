package repository

import (
	"context"
	"encoding/json"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderQueueRepository struct {
	channel *amqp.Channel
	queue   string
}

func NewOrderQueueRepository(channel *amqp.Channel, queue string) *OrderQueueRepository {
	return &OrderQueueRepository{
		channel: channel,
		queue:   queue,
	}
}

// PublishOrder sends an order to the RabbitMQ queue
func (r *OrderQueueRepository) PublishOrder(ctx context.Context, order entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
