//+build unit

package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/magiconair/properties/assert"
)

func TestNewConsumer(t *testing.T) {

}

func TestNewProducer(t *testing.T) {

}

func TestBuildConsumerConfig(t *testing.T) {
	//Act
	cfg := buildConsumerConfig()

	//Arrange
	assert.Equal(t, cfg.Version, sarama.V1_0_0_0)
	assert.Equal(t, cfg.Consumer.Offsets.Initial, sarama.OffsetOldest)
	assert.Equal(t, cfg.Consumer.Return.Errors, true)
	assert.Equal(t, cfg.Group.Return.Notifications, true)
	assert.Equal(t, cfg.Consumer.Fetch.Max, int32(1024*1024*2))
	assert.Equal(t, cfg.Consumer.Fetch.Default, int32(1024*512))
	assert.Equal(t, cfg.Consumer.Fetch.Min, int32(1024*10))
}

func TestBuildProducerConfig_WithMurmurHasher(t *testing.T) {

}

func TestBuildProducerConfig_WithFNVHasher(t *testing.T) {

}
