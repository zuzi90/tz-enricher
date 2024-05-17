package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"time"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	log      *logrus.Entry
}

func NewProducer(brokers []string, topic string, log *logrus.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 3
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		log:      log.WithField("module", "producer"),
	}, nil
}

func (p *Producer) SendMessage(msg []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := p.producer.SendMessage(message)
	if err != nil {
		p.log.Warnf("failed to send message: %v", err)
		return err
	}

	return err
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		p.log.Warnf("failed to close producer: %v", err)
		return err
	}

	return nil
}
