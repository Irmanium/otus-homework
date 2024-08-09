package rabbitrepo

import (
	"context"
	"encoding/json"

	"otus-homework/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "post live feed"
)

type Repo struct {
	conn *amqp.Connection
}

func New(conn *amqp.Connection) (*Repo, error) {
	r := &Repo{
		conn: conn,
	}

	err := r.declareExchange(conn)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repo) declareExchange(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *Repo) SendPost(ctx context.Context, post domain.Post) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		exchangeName,
		post.UserID,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (r *Repo) GetFeed(friendIDs []string, cancel <-chan struct{}) (<-chan []byte, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, friendID := range friendIDs {
		err = ch.QueueBind(
			q.Name,
			friendID,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	feed := make(chan []byte)
	go func() {
		defer func() {
			ch.Close()
			close(feed)
		}()

		select {
		case <-cancel:
			return
		case msg := <-msgs:
			feed <- msg.Body
		}
	}()

	return feed, nil
}
