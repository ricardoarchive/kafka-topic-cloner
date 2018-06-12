# Kafka topic cloner
`Kafka topic cloner` is a CLI that clones the content of a topic into another one.

The two topics have to coexist inside the same Kafka cluster. If that's not the case, you probably want to have a look at [Apache Mirrormaker](https://docs.confluent.io/current/multi-dc/mirrormaker.html).

The cloner supports two different hashers for partition assignment:
* `Murmur2`, which is the standard hasher used in the Java kafka community, including kafka scripts and kafka connect (default)
* `FNV-1a`, which is the standard hasher used in a part of the Go kafka community (_e.g. Sarama producers_)

The CLI was written in Go using [spf13/cobra](https://github.com/spf13/cobra).

## Installation

### Standard

It is strongly recommended that you use a released version. You can find the released binaries [here](https://github.com/ricardo-ch/kafka-topic-cloner/releases).

```sh
wget https://github.com/ricardo-ch/kafka-topic-cloner/releases/download/v0.5.0/kafka-topic-cloner
```

This project has binaries for Linux(386), and since the 0.5.0 for Windows. If your platform is not supported, you can build the project manually, or open an issue asking us to add your platform to the supported ones. If you used an official release, there is no extra step required before you can start cloning!

### Building from the sources

`Note: the project was created using Go 1.10.1, so you will need to have it installed before proceeding with the next steps.`

If you want to dig inside the code, or build an executable for your own platform, you can download the source code with `go get`.

```sh
go get -u github.com/ricardo-ch/kafka-topic-cloner
```

You then need to download the dependencies of the project. We are using [dep](https://github.com/golang/dep) as our dependency manager, so you will need to install it (instructions available on their own github repository).
Once `dep` is installed, retrieve the source files of the dependencies:
```sh
dep ensure --update
```

You can then build an executable for your own using `go build`. If you want to build for another platform, you will need to set the `GOOS` and `GOARCH` environment variables to match the target system specifications.

## Examples

### Standard use

A standard use of the `Kafka topic cloner` would look like this:

```sh
kafka-topic-cloner --from-brokers localhost:9092 --from foo --to bar
```

This is going to consume every event from the `foo` topic and produce them inside the `bar` topic.

If you choose the same hasher that was used to populate the source topic, and that you have the same number of partitions in the source and target topics, the cloned topic will be an exact replica of the original one. If some events were mistakenly placed on the wrong partition (_e.g. by manually producing them_), the cloning would place them back on the right one.

By default, `Kafka topic cloner` will use `Murmur2` as partitioning hasher. It is the algorithm that is used by the kafka scripts, kafka connect, and the majority of the Java kafka libraries (including the Stream API). If you prefer to use `FNV-1a` instead, which is the hasher implemented in [Sarama](https://github.com/Shopify/sarama) (Golang's most popular kafka library), you can use the hasher parameter:

```sh
kafka-topic-cloner --brokers localhost:9092 --from foo --to bar --hasher FNV-1a
```

If you would like to see another hasher implemented, feel free to open an issue about this!

### Cross-cluster cloning

You can clone a topic from a kafka cluster to a different one, by specifying the `--to-cluster` parameter:
```sh
kafka-topic-cloner --from-brokers localhost:9092 --to-brokers remote-cluster:9092 --from foo --to bar
```

### Timing out

Technically, a Kafka topic has no definite end, but it is nice to know when the application is done cloning every available event in the source topic. To do so, `Kafka topic cloner` comes with a built-in timeout that will close the application when it was unable to clone any event for a certain amount of time. The default timeout delay is 10 seconds. You can override this delay by using the `timeout` parameter:

```sh
kafka-topic-cloner --brokers localhost:9092 --from foo --to bar --timeout 5000
```

As of today, there is no way to completely disable the timeout (to be implemented soon).

### Loop-cloning

Loop-cloning, or same-topic cloning, is the action of cloning a topic into itself. Since it creates a continuous flow of new events inside the source topic, the cloning will never end and quickly multiply the number of events.
Since this can be quite a dangerous action if done unintentionally, this action it protected by the --loop parameter.
When loop-cloning, you should not specify the target topic, and the source topic will be used as target:
```sh
kafka-topic-cloner --from-brokers localhost:9092 --from foo --loop
```

## Parameters

You can find the complete list of parameters below:

Argument        | Shorthand | Description
-----------     | --------- | -----------
from-brokers    | F         | Semicolon-separated list of the source kafka brokers
from            | f         | Source topic's name
to-brokers      | T         | Semicolon-separated list of the target kafka brokers, specify only for cross-clusters cloning
to              | t         | Destination topic's name
timeout         | o         | consumer timeout is ms (defaults to 10000)
hasher          | H         | name of the hasher to use for partitioning, possible values: murmur2 (default), FNV-1a
loop            | L         | allow loop-cloning
verbose         | v         | verbose mode (defaults to false)
help            | h         | displays the CLI's help

## Disclaimer

`Kafka topic cloner` can not be considered production-ready since it does not come with any test at the moment. Use it at your own risk! :)

## Contributing

Contributions are always welcome and appreciated, the maintainers have a look at the new issues and consider every pull request.