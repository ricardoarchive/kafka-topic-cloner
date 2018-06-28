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

type parameters struct {
	verbose     bool
	loop        bool
	fromBrokers string
	toBrokers   string
	fromTopic   string
	toTopic     string
	hasher      string
	timeout     int
}

var (
	params          parameters
	consumerGroup   = "kafka-topic-cloner"
	possibleHashers = []string{"murmur2", "FNV-1a"}

	errMissingSourceTopic    = errors.New("source topic must be set")
	errMissingTargetTopic    = errors.New("target topic must be set")
	errLoopCloningWithTarget = errors.New("do not specify target topic when loop-cloning")
	errLoopRequired          = errors.New("cannot clone into the same topic without using --loop")
	errMissingSourceBrokers  = errors.New("source brokers must be set")
	errSourceBrokersIsTarget = errors.New("source and target brokers are identical")
	errUnknownHasher         = errors.New("unknown hasher, see help for possible value")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafka-topic-cloner --from-brokers [url] --from [source] --to [target]",
	Short: "Small cli used to clone the content of a topic into another one",
	Long: `
	Kafka topic cloner consumes all the events stored in the source topic, and produces them in the target topic.
	The events will see their partitions re-assigned in the process, based on the hasher that you can specify.

	Cloning between two different clusters can be achieved by using the --to-cluster flag.

	Same-topic cloning (also called loop-cloning) is protected by the --loop flag. In this case, the source topic (--from) will be used as both source and target.
	This can be a risky operation since it will multiply the messages in the source topic until manual interruption, use with caution!
	`,
	Run: Clone,
}

//Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//Load the cobra flags
	rootCmd.PersistentFlags().BoolVarP(&params.verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&params.loop, "loop", "L", false, "loop mode (clone into the source topic)")
	rootCmd.PersistentFlags().StringVarP(&params.fromBrokers, "from-brokers", "F", "", "address of the source kafka brokers, semicolon-separated")
	rootCmd.PersistentFlags().StringVarP(&params.toBrokers, "to-brokers", "T", "", "address of the target kafka brokers, semicolon-separated (specify only if different from the source brokers)")
	rootCmd.PersistentFlags().StringVarP(&params.fromTopic, "from", "f", "", "source topic")
	rootCmd.PersistentFlags().StringVarP(&params.toTopic, "to", "t", "", "target topic")
	rootCmd.PersistentFlags().StringVarP(&params.hasher, "hasher", "H", "murmur2", "partitioning hasher (possible values: murmur2, FNV-1a")
	rootCmd.PersistentFlags().IntVarP(&params.timeout, "timeout", "o", 10000, "delay (ms) before exiting after the last message has been cloned")

	rootCmd.MarkPersistentFlagRequired("from-brokers")
	rootCmd.MarkPersistentFlagRequired("from")
}

//Clone handles the consuming / producing process
func Clone(cmd *cobra.Command, args []string) {

	if err := params.validate(); err != nil {
		log.Print(err)
		return
	}

	fromBrokers, toBrokers := getBrokers()

	consumer := kafka.NewConsumer(params.fromTopic, fromBrokers, consumerGroup)
	if params.verbose {
		log.Printf("consumer (group: %s) initialized on %s/%s", consumerGroup, fromBrokers, params.fromTopic)
	}

	producer := kafka.NewProducer(toBrokers, params.hasher)
	if params.verbose {
		log.Printf("producer initialized on %s/%s, hasher: %s", toBrokers, params.toTopic, params.hasher)
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
			if params.verbose {
				log.Print(fmt.Sprintf("message consumed at partition %v, offset %v", msgC.Partition, msgC.Offset))
			}
			if ok {
				msgP := &sarama.ProducerMessage{
					Topic: params.toTopic,
				}
				if msgC.Value != nil {
					msgP.Value = sarama.ByteEncoder(msgC.Value)
				}
				if msgC.Key != nil {
					msgP.Key = sarama.ByteEncoder(msgC.Key)
				}
				producer.Input() <- msgP
				if params.verbose {
					log.Print("message produced")
				}
			}

		case <-signals:
			log.Print("terminating application")
			break Loop

		case <-time.After(time.Duration(params.timeout) * time.Millisecond):
			log.Print("timeout - end of cloning")
			break Loop
		}
	}
}

func (p parameters) validate() error {
	switch true {

	case p.fromTopic == "":
		return errMissingSourceTopic

	case p.toTopic == "" && !p.loop:
		return errMissingTargetTopic

	case p.toTopic != "" && p.loop:
		return errLoopCloningWithTarget

	case p.fromTopic == p.toTopic && !p.loop && p.toBrokers == "":
		return errLoopRequired

	case p.fromBrokers == "":
		return errMissingSourceBrokers

	case p.fromBrokers == p.toBrokers:
		return errSourceBrokersIsTarget

	case !contains(possibleHashers, p.hasher):
		return errUnknownHasher

	}
	return nil
}

func getBrokers() (from, to []string) {
	from = strings.Split(params.fromBrokers, ";")

	if params.toBrokers != "" {
		to = strings.Split(params.toBrokers, ";")
	} else {
		to = from
	}
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
