package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

//NewConsumer configures and returns a cluster-consumer
func NewConsumer(from string, brokers []string, consumerGroup string) *cluster.Consumer {

	cfg := buildConsumerConfig()

	consumer, err := cluster.NewConsumer(brokers, consumerGroup, []string{from}, cfg)
	if err != nil {
		log.Fatal(err)
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

//NewProducer configures and returns an async producer
func NewProducer(brokers []string, hasher string) sarama.AsyncProducer {

	cfg := buildProducerConfig(hasher)

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to produce message: %+v\n", err)
		}
	}()

	go func() {
		for range producer.Successes() {
			//Safety first! :D
			//If return.Successes is set to true, not listening to this topic will
			//prevent the application from cloning after a certain amount of events.
		}
	}()

	return producer
}

func buildConsumerConfig() *cluster.Config {
	cfg := cluster.NewConfig()

	cfg.Version = sarama.V1_0_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	cfg.Consumer.Return.Errors = true
	cfg.Group.Return.Notifications = true

	// Without these, cloning a high-volume topic will fail
	cfg.Consumer.Fetch.Max = 1024 * 1024 * 2 //2 Mo
	cfg.Consumer.Fetch.Default = 1024 * 512
	cfg.Consumer.Fetch.Min = 1024 * 10

	return cfg
}

func buildProducerConfig(hasher string) *sarama.Config {

	cfg := sarama.NewConfig()

	//Has to be greater than 1_0_0_0 to send producer timestamps
	cfg.Version = sarama.V1_0_0_0
	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForLocal

	// Since producer messages size limits are best-effort only, we set the sarama absolute hard limit ti 1M
	sarama.MaxRequestSize = 1024 * 1024

	// Increasing this value will greatly increase the cloning speed.
	// However, with MaxOpenRequests > 1, the order of the cloned messages is not guaranteed.
	cfg.Net.MaxOpenRequests = 1

	// Without this, cloning a high-volume topic will fail
	cfg.Producer.Flush.Frequency = 100 * time.Millisecond

	if hasher == "murmur2" {
		cfg.Producer.Partitioner = sarama.NewCustomHashPartitioner(MurmurHasher)
	}

	return cfg
}
