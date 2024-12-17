package queuesample

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/HoangNguyen689/interface-practice/pkg/queue"
	"github.com/HoangNguyen689/interface-practice/pkg/queue/queueredis"
	"github.com/HoangNguyen689/interface-practice/pkg/queue/queuesqs"
)

const (
	standardQueueName = "test-standard-queue"
	fifoQueueName     = "test-fifo-queue.fifo"
	redisQueueName    = "test-redis-queue"
)

type queuesample struct{}

func NewCommand() *cobra.Command {
	q := &queuesample{}

	cmd := &cobra.Command{
		Use:   "queue-sample",
		Short: "Test queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := q.run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (q *queuesample) run() error {
	var ctx = context.Background()

	if err := testStandardQueue(ctx); err != nil {
		panic(err)
	}

	if err := testFifoQueue(ctx); err != nil {
		panic(err)
	}

	if err := testRedisQueue(ctx); err != nil {
		panic(err)
	}

	return nil
}

func testStandardQueue(ctx context.Context) error {
	fmt.Println("Stated testing the standard queue.")
	msgList := []string{
		"msg-1",
		"msg-2",
		"msg-3",
	}

	// init
	client, err := queuesqs.NewClient(ctx)
	if err != nil {
		return err
	}

	standardQueue, err := client.Queue(ctx, standardQueueName)
	if err != nil {
		return err
	}

	// send messages in order
	for _, msg := range msgList {
		if err := standardQueue.SendMessage(ctx, msg); err != nil {
			return err
		}
	}

	// receive messages. test 6 times. it should receive all messages in 3 times.
	// because the messages is distributed, and it change the subnet of server each time to cover all.
	fmt.Println("Received standard queue messages.")
	for i := 0; i < 6; i++ {
		time.Sleep(1 * time.Second)
		res, err := standardQueue.ReceiveMessages(ctx,
			queue.WithMaxNumberOfMessages(queue.MaxNumberOfMessages),
		)
		if err != nil {
			return err
		}

		fmt.Printf("Time %2d: received %2d messages\n", i+1, len(res))
		if len(res) > 0 {
			for i, msg := range res {
				fmt.Printf("%2d: %s\n", i+1, msg.Body)
			}

			// Delete messages
			for _, msg := range res {
				if err := standardQueue.DeleteMessage(ctx, msg.AckToken); err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("Tested the standard queue!")

	return nil
}

func testFifoQueue(ctx context.Context) error {
	fmt.Println("Started testing the fifo queue.")
	fifoMsgList := []string{
		"g-1-msg-1",
		"g-1-msg-2",
		"g-1-msg-3",
		"g-1-msg-4",
		"g-1-msg-5",
		"g-1-msg-6",
		"g-1-msg-7",
		"g-1-msg-8",
		"g-2-msg-1",
		"g-2-msg-2",
		"g-2-msg-3",
		"g-2-msg-4",
		"g-2-msg-5",
		"g-2-msg-6",
		"g-1-msg-9",
	}

	// init
	fifoClient, err := queuesqs.NewClient(ctx)
	if err != nil {
		return err
	}

	fifoQueue, err := fifoClient.Queue(ctx, fifoQueueName)
	if err != nil {
		return err
	}

	// send messages in order
	for _, msg := range fifoMsgList {
		group := "g-1"
		if strings.HasPrefix(msg, "g-2") {
			group = "g-2"
		}

		dedupID := uuid.New().String()

		if err := fifoQueue.SendMessage(ctx, msg,
			queue.WithGroupID(group),
			queue.WithDeDuplicationID(dedupID),
		); err != nil {
			return err
		}
	}

	// receive messages
	res, err := fifoQueue.ReceiveMessages(ctx,
		queue.WithMaxNumberOfMessages(queue.MaxNumberOfMessages),
	)
	if err != nil {
		return err
	}

	fmt.Println("Received fifo queue messages:")
	for i, msg := range res {
		fmt.Printf("%2d: %s\n", i+1, msg.Body)
	}

	// at this time, it should receive all messages from group 1
	// because it consumes as many messages as possible from one group first.
	// we can't received any messages from g-2-msg-2 anymore
	// because the g-2-msg-1 is not deleted yet.
	// test for 6 times
	for i := 0; i < 6; i++ {
		time.Sleep(1 * time.Second)

		res, err := fifoQueue.ReceiveMessages(ctx,
			queue.WithMaxNumberOfMessages(queue.MaxNumberOfMessages),
		)
		if err != nil {
			return err
		}

		fmt.Printf("Time %2d: received %2d messages\n", i+1, len(res))
		if len(res) > 0 {
			fmt.Printf("Can retrieve messages from fifo queue!!!\n")
			for i, msg := range res {
				fmt.Printf("%2d: %s", i+1, msg.Body)
			}
		}
	}

	// delete all the messages except "g-1-msg-6", "g-1-msg-2", put it visible again
	for _, msg := range res {
		if msg.Body == "g-1-msg-6" {
			fmt.Println("Got the g-1-msg-6 message, put it visible again")
			fifoQueue.ChangeMessageVisibility(ctx, msg.AckToken, 0)

			continue
		}

		if msg.Body == "g-1-msg-2" {
			fmt.Println("Got the g-1-msg-2 message, put it visible again with longer visibility timeout")
			fifoQueue.ChangeMessageVisibility(ctx, msg.AckToken, 5)

			continue
		}

		if err := fifoQueue.DeleteMessage(ctx, msg.AckToken); err != nil {
			return err
		}
	}

	// receive messages again. at this point, it will receive all the group 2 messages
	// and the g-1-msg-2, g-1-msg-6 message at last.
	// g-1-msg-2 still go before g-1-msg-6.
	for i := 0; i < 6; i++ {
		time.Sleep(1 * time.Second)

		res, err := fifoQueue.ReceiveMessages(ctx,
			queue.WithMaxNumberOfMessages(queue.MaxNumberOfMessages),
		)

		if err != nil {
			return err
		}

		fmt.Printf("Time %2d: received %2d messages\n", i+1, len(res))
		if len(res) > 0 {
			fmt.Printf("Can retrieve messages from fifo queue!!!\n")
			for i, msg := range res {
				fmt.Printf("%2d: %s\n", i+1, msg.Body)
			}

			for _, msg := range res {
				if err := fifoQueue.DeleteMessage(ctx, msg.AckToken); err != nil {
					return err
				}
			}
		}

	}
	fmt.Println("Tested the standard queue!")

	return nil
}

func testRedisQueue(ctx context.Context) error {
	fmt.Println("Stated testing the redis queue.")
	msgList := []string{
		"msg-1",
		"msg-2",
		"msg-3",
	}

	// init
	client, err := queueredis.NewClient(ctx, "localhost:6379")
	if err != nil {
		return err
	}

	redisQueue, err := client.Queue(ctx, redisQueueName)
	if err != nil {
		return err
	}

	// send messages in order
	for _, msg := range msgList {
		if err := redisQueue.SendMessage(ctx, msg); err != nil {
			return err
		}
	}

	// receive messages. test 6 times. it should receive all messages in 3 times.
	// and the last 3 times should be empty
	fmt.Println("Received redis queue messages.")
	for i := 0; i < 6; i++ {
		time.Sleep(1 * time.Second)
		res, err := redisQueue.ReceiveMessages(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Time %2d: received %2d messages\n", i+1, len(res))
		if len(res) > 0 {
			for i, msg := range res {
				fmt.Printf("%2d: %s\n", i+1, msg.Body)
			}
		}
	}

	fmt.Println("Tested the redis queue!")

	return nil
}
