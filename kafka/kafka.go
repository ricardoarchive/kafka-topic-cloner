package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

//NewConsumer returns a new cluster consumer
func NewConsumer(from string, brokers []string, consumerGroup string) *cluster.Consumer {
	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Group.Return.Notifications = true

	consumer, err := cluster.NewConsumer(brokers, consumerGroup, []string{from}, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	return consumer
}

//NewProducer returns a new producer
func NewProducer(brokers []string, defaultHasher bool) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	if defaultHasher {
		cfg.Producer.Partitioner = sarama.NewCustomHashPartitioner(MurmurHasher)
	}

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		panic(err)
	}

	return producer
}