package queue

import "context"

// This might be a little confusing.
// Just image that Client have method GetQueue, but the verb Get is unnecessary.
type Client interface {
	Queue(ctx context.Context, name string) (Queue, error)
}

type Queue interface {
	SendMessage(ctx context.Context, body string, opts ...SendMessagesOption) error
	ReceiveMessages(ctx context.Context, opts ...ReceiveMessagesOption) ([]Message, error)
	DeleteMessage(ctx context.Context, ackToken string) error

	ChangeMessageVisibility(ctx context.Context, ackToken string, visibilityTimeout int32) error
}
