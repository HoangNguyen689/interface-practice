package queuesqs

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/HoangNguyen689/interface-practice/pkg/queue"
)

type Client struct {
	*sqs.Client
}

type Queue struct {
	cli  *Client
	name string
	url  string
}

func NewClient(ctx context.Context) (*Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		sqs.NewFromConfig(awsCfg),
	}, nil
}

func (c *Client) Queue(ctx context.Context, name string) (queue.Queue, error) {
	res, err := c.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: &name})
	if err != nil {
		return nil, err
	}

	return Queue{
		cli:  c,
		name: name,
		url:  *res.QueueUrl,
	}, nil
}

func (q Queue) SendMessage(ctx context.Context, body string, opts ...queue.SendMessagesOption) error {
	attr := queue.SendMessageAttribute{}

	for _, opt := range opts {
		opt(&attr)
	}

	in := sqs.SendMessageInput{
		MessageBody:            &body,
		QueueUrl:               &q.url,
		MessageGroupId:         attr.GroupID,
		MessageDeduplicationId: attr.DeDuplicationID,
	}

	if _, err := q.cli.SendMessage(ctx, &in); err != nil {
		return err
	}

	return nil
}

func (q Queue) ReceiveMessages(ctx context.Context, opts ...queue.ReceiveMessagesOption) ([]queue.Message, error) {
	attr := queue.ReceiveMessageAttribute{}

	for _, opt := range opts {
		opt(&attr)
	}

	in := sqs.ReceiveMessageInput{
		QueueUrl: &q.url,
		AttributeNames: []types.QueueAttributeName{ // must specify these attributes to get the value
			"MessageGroupId",
			"MessageDeduplicationId",
		},
		MaxNumberOfMessages: attr.MaxNumberOfMessages,
	}

	res, err := q.cli.ReceiveMessage(ctx, &in)
	if err != nil {
		return nil, err
	}

	msgList := make([]queue.Message, 0, len(res.Messages))
	for _, m := range res.Messages {
		msgList = append(msgList, queue.Message{
			ID:              *m.MessageId,
			GroupID:         m.Attributes["MessageGroupId"],
			DeDuplicationID: m.Attributes["MessageDeduplicationId"],
			AckToken:        *m.ReceiptHandle,
			Body:            *m.Body,
		})
	}

	return msgList, nil
}

func (q Queue) DeleteMessage(ctx context.Context, ID string) error {
	in := sqs.DeleteMessageInput{
		QueueUrl:      &q.url,
		ReceiptHandle: &ID,
	}

	if _, err := q.cli.DeleteMessage(ctx, &in); err != nil {
		return err
	}

	return nil
}

func (q Queue) ChangeMessageVisibility(ctx context.Context, ackToken string, visibilityTimeout int32) error {
	in := sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &q.url,
		ReceiptHandle:     &ackToken,
		VisibilityTimeout: visibilityTimeout,
	}

	if _, err := q.cli.ChangeMessageVisibility(ctx, &in); err != nil {
		return err
	}

	return nil
}
