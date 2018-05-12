// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ricardo-ch/kafka-topic-cloner/kafka"
	"github.com/spf13/cobra"
)

var (
	verbose       bool
	defaultHasher bool
	url           string
	from          string
	to            string
	timeout       int
	consumerGroup = "kafka-topic-cloner"
)

const appName = "kafka-topic-cloner"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafka-topic-cloner",
	Short: "small cli used to clone the content of a topic into another one",
	Long: `consumes all the events stored in the input topic,
	and produces them in the output topic.
	The two topics have to co-exist inside the same cluster.`,
	Run: Clone,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&defaultHasher, "default-hasher", "d", false, "use the default sarama hasher for partitioning instead of murmur2")
	rootCmd.PersistentFlags().StringVarP(&url, "brokers", "b", "http://localhost:9092", "semicolon-separated Kafka brokers URLs")
	rootCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "input topic")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "output topic")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "o", 10000, "delay before timing out")
}

//Clone ...
func Clone(cmd *cobra.Command, args []string) {

	brokers := strings.Split(url, ";")
	consumer := kafka.NewConsumer(from, brokers, consumerGroup)
	if verbose {
		log.Printf("consumer %s initialized on %s/%s", consumerGroup, brokers, from)
	}
	producer := kafka.NewProducer(brokers, defaultHasher)
	if verbose {
		log.Printf("producer initialized on %s/%s", brokers, to)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	//Capture interrupt signal to stop the application
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msgC, ok := <-consumer.Messages():
			if verbose {
				log.Print("message consumed")
			}
			if ok {
				msgP := &sarama.ProducerMessage{
					Topic: to,
					Key:   sarama.ByteEncoder(msgC.Key),
					Value: sarama.ByteEncoder(msgC.Value),
				}
				_, _, err := producer.SendMessage(msgP)
				if err != nil {
					log.Printf("FAILED to send message %s\n", err)
				}
				if verbose {
					log.Print("message produced")
				}
			}
		case <-signals:
			log.Print("terminating application")
			return
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			log.Print("timeout - end of cloning")
			return
		}
	}
}
