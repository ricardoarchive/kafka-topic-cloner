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
	"errors"
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
	force         bool
	brokers       string
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
	rootCmd.PersistentFlags().BoolVarP(&defaultHasher, "force", "F", false, "force cloning into the source topic")
	rootCmd.PersistentFlags().StringVarP(&brokers, "brokers", "b", "", "semicolon-separated Kafka brokers URLs")
	rootCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "source topic")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "target topic")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "o", 10000, "delay before timing out")

	rootCmd.MarkFlagRequired("brokers")
	rootCmd.MarkFlagRequired("from")
	rootCmd.MarkFlagRequired("to")
}

//Clone handles the consuming / producing process
func Clone(cmd *cobra.Command, args []string) {

	if err := validateParameters(); err != nil {
		log.Print(err)
		return
	}

	b := strings.Split(brokers, ";")
	consumer := kafka.NewConsumer(from, b, consumerGroup)
	if verbose {
		log.Printf("consumer (group: %s) initialized on %s/%s", consumerGroup, b, from)
	}
	producer := kafka.NewProducer(b, defaultHasher)
	if verbose {
		log.Printf("producer initialized on %s/%s, default hasher: %t", b, to, defaultHasher)
	}

	//Try to gracefully shutdown
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//Capture interrupt and kill signal to stop the application
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	//Cloning loop
Loop:
	for {
		select {

		case msgC, ok := <-consumer.Messages():
			if verbose {
				log.Print(fmt.Sprintf("message consumed at partition %v, offset %v", msgC.Partition, msgC.Offset))
			}
			if ok {
				msgP := &sarama.ProducerMessage{
					Topic: to,
					Key:   sarama.ByteEncoder(msgC.Key),
					Value: sarama.ByteEncoder(msgC.Value),
				}
				producer.Input() <- msgP
				if verbose {
					log.Print("message produced")
				}
			}

		case <-signals:
			log.Print("terminating application")
			break Loop

		case <-time.After(time.Duration(timeout) * time.Millisecond):
			log.Print("timeout - end of cloning")
			break Loop
		}
	}
}

func validateParameters() error {
	switch true {
	case brokers == "":
		return errors.New("brokers must be set")
	case from == "":
		return errors.New("source topic must be set")
	case to == "":
		return errors.New("target topic must be set")
	case from == to && !force:
		return errors.New("same-topic cloning is disabled, use --force to ignore")
	}

	return nil
}
