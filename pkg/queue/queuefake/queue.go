package queuefake

import (
	"context"

	"github.com/HoangNguyen689/interface-practice/pkg/queue"
)

type Client struct{}

type Queue struct{}

func NewClient(ctx context.Context) (*Client, error) {
	return &Client{}, nil
}

func (c *Client) Queue(ctx context.Context, name string) (queue.Queue, error) {
	return Queue{}, nil
}

func (q Queue) SendMessage(ctx context.Context, body string, opts ...queue.SendMessagesOption) error {
	return nil
}

func (q Queue) ReceiveMessages(ctx context.Context, opts ...queue.ReceiveMessagesOption) ([]queue.Message, error) {
	return []queue.Message{
		{
			ID:   "id-1",
			Body: `{"hello": "world"}`,
		},
	}, nil
}

func (q Queue) DeleteMessage(ctx context.Context, ackToken string) error {
	return nil
}

func (q Queue) ChangeMessageVisibility(ctx context.Context, ackToken string, visibilityTimeout int32) error {
	return nil
}
