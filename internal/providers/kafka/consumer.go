package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type messageHandler interface {
	Handle(ctx context.Context, fio []byte) error
}

type Consumer struct {
	messageHandler messageHandler
	consumer       sarama.Consumer
	log            *logrus.Entry
	poolCh         chan func()
	workersCount   int
	kafkaTopic     string
	brokers        []string
}

func NewConsumer(brokers []string, wCount int, kafkaTopic string, log *logrus.Logger, messageHandler messageHandler) *Consumer {
	c := Consumer{
		messageHandler: messageHandler,
		workersCount:   wCount,
		kafkaTopic:     kafkaTopic,
		brokers:        brokers,
		log:            log.WithField("module", "consumer"),
	}

	c.poolCh = make(chan func(), c.workersCount)
	for i := 0; i < c.workersCount; i++ {
		go func() {
			for f := range c.poolCh {
				f()
			}
		}()
	}

	c.log.Infof("consumer is ready to consume messages from topic %s", c.kafkaTopic)

	return &c
}

func (c *Consumer) Run(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Consumer.IsolationLevel = sarama.ReadCommitted
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	consumer, err := sarama.NewConsumer(c.brokers, config)
	if err != nil {
		return err
	}

	c.consumer = consumer
	partitions, err := c.consumer.Partitions(c.kafkaTopic)
	if err != nil {
		c.log.Warnf("partitions are not available %v", err)

		return err
	}

	for _, partition := range partitions {
		consumePartition, err := c.consumer.ConsumePartition(c.kafkaTopic, partition, sarama.OffsetNewest)
		if err != nil {
			c.log.Warnf("partitions are not available %v", err)

			return err
		}
		go func(consumePartition sarama.PartitionConsumer) error {
			defer consumePartition.Close()

			for message := range consumePartition.Messages() {

				select {
				case <-ctx.Done():
					return nil
				default:
				}

				c.poolCh <- func() {

					if err := c.messageHandler.Handle(ctx, message.Value); err != nil {
						c.log.Warnf("handling message: %v: %v", string(message.Value), err)
					}

				}
			}
			return nil

		}(consumePartition)

	}

	return nil
}
