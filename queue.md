# AWS SQS

## Message

An Amazon SQS message has three basic states:
1. Sent to a queue by a producer.
2. Received from the queue by a consumer.
3. Deleted from the queue.

State between 2 and 3 considered as **in-flight**. Can monitor in-flight message number in Cloudwatch.

Message retention
- a message is retained for 4 days.
- the minimum is 60 seconds (1 minute).
- the maximum is 1,209,600 seconds (14 days).

Message attributes
- Can include structured metadata (such as timestamps, geospatial data, signatures, and identifiers) with messages.


## Standard queue

- at-least-one message (message can be duplicated, out of order).
- design the application **idempotent**
- examples that fit
  - Decoupling live user requests from intensive background work.
  - Allocating tasks (i.e. validation) to multiple worker nodes.
  - Batching messages for future processing.

## FIFO queue (first in first out)

- exactly-one message.
- support group messages. allow multiple ordered message groups within a single queue.
- examples that fit
  - E-commerce order management system where order is critical
  - Processing user-entered inputs in the order entered
  - Communications and networking – Sending and receiving data and information in the same order
  - Online ticketing system
  - Educational institutes

### Term

**Group ID**

- the messages are ordered inside a group. the different group messages may be out of order.
- max 128 characters
- Valid values: alphanumeric characters and punctuation
(!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~)

**Deduplication ID**

- multiple messages with the same deduplication ID are sent within a 5 minute
deduplication interval, they are treated as duplicates, and only one copy is delivered.
- the deduplication interval is fixed for 5 mins.
- specify explicitly or enable the content-based deduplication


### Receive messages

- cannot explicitly request to receive messages from a specific message group ID.
- when have many groups:
  - attempts to return as many messages as possible with the same message group ID in a single call.
  - allows other consumers to process messages from different message group IDs concurrently.
  - may receive multiple messages from the same message group ID in one batch (up to 10)
  - can't receive additional messages from the same message group ID in subsequent requests until:
    - the currently received messages are deleted, or
    - they become visible again (for example, after the visibility timeout expires)
  - if fewer than 10 messages are available for the same message group ID, may include messages from other message group IDs in the same batch, but each group retains order.

### Retry

- producer can retry SendMessage freely using the same message deduplication ID before the
deduplication interval expires.
- consumer can retry ReceiveMessage freely using the same receive request attempt ID before the visibility timeout expires.

### Visibility timout

- when a message is retrieved but not deleted, it remains invisible until the visibility timeout
expires.
- no additional messages from the same message group ID are returned until the first message
is deleted or becomes visible again.
- the default visibility timeout for a queue is 30 seconds,
- the maximum visibility timeout is 12 hours, fixed from when the message is first received.
- can extend the visibility timeout but not exceed the hard limit.

### MaxNumberOfMessages

- the maximum number of messages to return. never returns more messages than
this value (however, fewer messages might be returned).
- from 1 to 10
- default is 1
- may set to 1 for simplicity

## Dead-letter queue

- dead-letter queues (DLQs), which source queues can target for messages that are not processed successfully.
- maxReceiveCount set up for the DLQs
  - if the message is received with times bigger than maxReceiveCount, it will be moved to DLQs
  - notice about the visibility timeout
  - don't recommend set to 1 because the queue is distributed and in case at-least-one, there're not any gaurantee that the message can be received by ReceiveMessage.
- redrive to move message back to original queue.



## Delay queue

- postpone the delivery of new messages to consumers for a number of seconds.
- the default is 0 second.
- the maximum is 15 minutes

## Short and long polling

- for receiving messages from a queue
- short polling (default): queries a subset of servers, receive messages immediately
- long polling: queries all servers for messages, sending a response once at
least one message is available, up to the specified maximum. An empty response is sent only if the polling wait time expires.
  - the maximum long polling wait time is 20 seconds.

## Message timer

- delay the time consumers receive message
- not apply for individual message on FIFO
- the maximum is 15 minutes

## Security

### Attribute-based access control (ABAC)

- access control based on attributes
  - authorization by tags and aliases.
- default to AWS SQS managed server-side encryption (SSE) option
  - can create custom managed server-side encryption that uses
SQS-managed encryption keys to protect sensitive data sent over message queues.

## Note

- if have problems about throughput, consider `High throughput for FIFO queue`

## Ref

- https://docs.aws.amazon.com/pdfs/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-dg.pdf
- https://docs.aws.amazon.com/pdfs/AWSSimpleQueueService/latest/APIReference/sqs-api.pdf
- https://aws.amazon.com/blogs/compute/solving-complex-ordering-challenges-with-amazon-sqs-fifo-queues/