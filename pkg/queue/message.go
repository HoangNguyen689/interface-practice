package queue

type Message struct {
	ID              string
	GroupID         string
	DeDuplicationID string
	AckToken        string
	Body            string
}
