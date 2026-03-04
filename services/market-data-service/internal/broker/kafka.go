package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type PriceEvent struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	UpdatedAt int64   `json:"updated_at"`
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}

func (k *KafkaProducer) Publish(symbol string, price float64) error {
	event := PriceEvent{
		Symbol:    symbol,
		Price:     price,
		UpdatedAt: time.Now().Unix(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(symbol),
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = k.producer.SendMessage(msg)
	return err
}
