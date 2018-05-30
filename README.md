# Kafka topic cloner
`Kafka topic cloner` is a CLI that clones the content of a topic into another one.

The two topics have to coexist inside the same Kafka cluster. If that's not the case, you probably want to have a look at [Apache Mirrormaker](https://docs.confluent.io/current/multi-dc/mirrormaker.html).

The cloner supports two different hashers for partition assignment:
* `Murmur2`, which is the standard hasher used in the Java kafka community, including kafka scripts and kafka connect (default)
* `FNV-1a`, which is the standard hasher used in the Go kafka community (_e.g. Sarama producers_)

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
Once installed, a simple `dep ensure` will retrieve the source files of the dependencies.

You can then build an executable for your own using `go build`. If you want to build for another platform, you will need to set the `GOOS` and `GOARCH` environment variables to match the target system specifications.

## Usage

A standard use of the `Kafka topic cloner` would look like this:

```sh
kafka-topic-cloner --brokers localhost:9092 --from foo --to bar
```

This is going to consume every message from the `foo` topic and produce them inside the `bar` topic.

If you choose the same hasher that was used to populate the source topic, and that you have the same number of partitions in the source and target topics, the cloned topic will be an exact replica of the original one. If some messages were mistakenly placed on the wrong partition (e.g. by manually producing them), the cloning would place them back on the right one.

By default, `Kafka topic cloner` will use `Murmur2` as hasher algorithm. It is the algorithm that is used by the kafka scripts, kafka connect, and the majority of the Java kafka libraries (including the Stream API). If you prefer to use `FNV-1a` instead, which is the hasher implemented in [Sarama](https://github.com/Shopify/sarama) (Golang's most popular kafka library), you can use the default-hasher parameter:

```sh
kafka-topic-cloner --brokers localhost:9092 --from foo --to bar --default-hasher
```

If you would like to see another hasher implemented, feel free to open an issue about this.

Technically, a Kafka topic has no definite end, but it is nice to know when the application is done cloning every available message in the source topic. To do so, `Kafka topic cloner` comes with a built-in timeout that will close the application when it was unable to clone any message for a certain amount of time. The default timeout delay is 10 seconds. You can override this delay by using the `timeout` parameter:

```sh
kafka-topic-cloner --brokers localhost:9092 --from foo --to bar --timeout 5000
```

As of today, there is no way to completely disable the timeout.

You can find the complete list of parameters below:

Argument        | Shorthand | Description
-----------     | --------- | -----------
brokers         | b         | Semicolon-separated list of kafka brokers `required`
from            | f         | Source topic's name `required`
to              | t         | Destination topic's name`required`
timeout         | o         | consumer timeout is ms (defaults to 10000)
default-hasher  | d         | use Sarama's FNV-1a hasher (defaults to false)
verbose         | v         | verbose mode (defaults to false)

## Disclaimer

`Kafka topic cloner` can not be considered production-ready since it does not come with any test at the moment. Use it at your own risk! :)

## Contributing

Contributions are always welcome and appreciated, the maintainers have a look at the new issues and consider every pull request.