package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/cyub/hyper/pkg/queue"
)

// Queue struct
type Queue struct {
	*queue.Base
	Options
	producer sarama.SyncProducer
	ctx      context.Context
	cancel   context.CancelFunc
	messages chan string
}

// Options struct
type Options struct {
	queue.Options
	Brokers      []string
	GroupID      string
	KafkaVersion string
}

var _ queue.Queuer = (*Queue)(nil)

// New return queue base kafka
func New(opts Options) *Queue {
	opts.Init()
	queue, err := newQueue(opts)
	if err != nil {
		panic(err)
	}
	return queue
}

func newQueue(opts Options) (*Queue, error) {
	q := &Queue{
		Options:  opts,
		messages: make(chan string),
	}

	q.Base = queue.NewBase(opts.Options)
	q.Base.Enqueue = q.doEnqueue
	producer, err := q.createProducer()
	if err != nil {
		return nil, err
	}
	q.producer = producer
	client, err := q.createConsumerGroup()
	if err != nil {
		return nil, err
	}
	q.joinConsumerCluster(client)
	return q, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (q *Queue) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (q *Queue) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (q *Queue) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if q.Debug {
			fmt.Printf("message claimed: topic[%s], partition[%d], timestamp[%v], offset[%d], body[%s]\n", message.Topic, message.Partition, message.Timestamp, message.Offset, string(message.Value))
		}
		q.messages <- string(message.Value)
		session.MarkMessage(message, "")
	}

	return nil
}

func (q *Queue) createProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	return sarama.NewSyncProducer(q.Brokers, config)
}

func (q *Queue) createConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(q.KafkaVersion)
	if err != nil {
		q.Logger.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return sarama.NewConsumerGroup(q.Brokers, q.GroupID, config)
}

func (q *Queue) joinConsumerCluster(client sarama.ConsumerGroup) {
	q.ctx, q.cancel = context.WithCancel(context.Background())
	go func() {
		for {
			if err := client.Consume(q.ctx, []string{q.Name}, q); err != nil {
				q.Logger.Errorf("join the consumer cluster of topic[%s] fail %s", q.Name, err.Error())
				q.Backoff()
			} else {
				q.BackoffReset()
			}
		}
	}()
}

func (q *Queue) doEnqueue(msg []byte) error {
	partition, offset, err := q.producer.SendMessage(&sarama.ProducerMessage{
		Topic: q.Name,
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		return err
	}

	if q.Debug {
		fmt.Printf("kafka message enqueue: topic[%s], partition[%d], offset[%d]\n", q.Name, partition, offset)
	}
	return nil
}

// Run use for queue run
func (q *Queue) Run() {
	go func() {
		for {
			data := <-q.messages
			job, err := q.Parse(data)
			if err != nil {
				q.Logger.Errorf("job[%s] parse failure %s", data, err.Error())
				continue
			}
			fn, err := q.GetConsumer(job.GetName())
			if err != nil {
				q.Logger.Error(err)
			}
			q.SemAcquire(1)
			go q.Do(job, fn)
		}
	}()
}
