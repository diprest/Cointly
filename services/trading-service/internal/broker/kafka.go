package broker

import (
	"encoding/json"
	"log"
	"trading-service/internal/models"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewKafkaConsumer(brokers []string, topic string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	c, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, topic: topic}, nil
}

func (k *KafkaConsumer) Subscribe(updates chan<- models.PriceUpdate) {
	partitions, err := k.consumer.Partitions(k.topic)
	if err != nil {
		log.Printf("Failed to get partitions: %v", err)
		return
	}

	for _, p := range partitions {
		pc, err := k.consumer.ConsumePartition(k.topic, p, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to consume partition %d: %v", p, err)
			continue
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var price models.PriceUpdate
				if err := json.Unmarshal(msg.Value, &price); err == nil {
					updates <- price
				}
			}
		}(pc)
	}
	log.Println("Kafka Consumer started listening...")
}
