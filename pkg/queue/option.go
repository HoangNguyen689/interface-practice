package queue

const (
	MaxNumberOfMessages = 10
)

type SendMessageAttribute struct {
	GroupID         *string
	DeDuplicationID *string
}

type SendMessagesOption func(a *SendMessageAttribute)

func WithGroupID(ID string) SendMessagesOption {
	return func(a *SendMessageAttribute) {
		a.GroupID = &ID
	}
}

func WithDeDuplicationID(ID string) SendMessagesOption {
	return func(a *SendMessageAttribute) {
		a.DeDuplicationID = &ID
	}
}

type ReceiveMessageAttribute struct {
	MaxNumberOfMessages int32
	VisibilityTimeout   int32 // for extend the visibility timeout
	WaitTimeSeconds     int32 // for long polling, short polling has this value = 0
}

type ReceiveMessagesOption func(a *ReceiveMessageAttribute)

func WithMaxNumberOfMessages(n int32) ReceiveMessagesOption {
	return func(a *ReceiveMessageAttribute) {
		a.MaxNumberOfMessages = n
	}
}

func WithWaitTimeSeconds(n int32) ReceiveMessagesOption {
	return func(a *ReceiveMessageAttribute) {
		a.WaitTimeSeconds = n
	}
}

func WithVisibilityTimeout(n int32) ReceiveMessagesOption {
	return func(a *ReceiveMessageAttribute) {
		a.VisibilityTimeout = n
	}
}
