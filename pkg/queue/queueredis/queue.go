package queueredis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/HoangNguyen689/interface-practice/pkg/queue"
)

type Client struct {
	*redis.Client
}

type Queue struct {
	cli *Client
	key string
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Client{
		r,
	}, nil
}

func (c *Client) Queue(ctx context.Context, key string) (queue.Queue, error) {
	return Queue{
		cli: c,
		key: key,
	}, nil
}

func (q Queue) SendMessage(ctx context.Context, body string, opts ...queue.SendMessagesOption) error {
	attr := queue.SendMessageAttribute{}

	for _, opt := range opts {
		opt(&attr)
	}

	// todo: options not being used

	if _, err := q.cli.LPush(ctx, q.key, body).Result(); err != nil {
		return err
	}

	return nil
}

func (q Queue) ReceiveMessages(ctx context.Context, opts ...queue.ReceiveMessagesOption) ([]queue.Message, error) {
	attr := queue.ReceiveMessageAttribute{}

	for _, opt := range opts {
		opt(&attr)
	}

	// todo: options not being used

	item, err := q.cli.RPop(ctx, q.key).Result()
	if err == redis.Nil {
		return []queue.Message{}, nil
	}
	if err != nil {
		return nil, err
	}

	return []queue.Message{
		{
			Body: item,
		},
	}, nil
}

func (q Queue) DeleteMessage(ctx context.Context, ID string) error {
	return nil
}

func (q Queue) ChangeMessageVisibility(ctx context.Context, ackToken string, visibilityTimeout int32) error {
	return nil
}
