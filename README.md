# Kafka topic cloner
`Kafka topic cloner` is a CLI that clones the content of a topic into another one.

The two topics have to coexist inside the same cluster. If that's not the case, you probably want to have a look at [Apache Mirrormaker](https://docs.confluent.io/current/multi-dc/mirrormaker.html).

The cloner supports two different hashers for partition assignment:
* `Murmur2`, which is the standard hasher used in the Java kafka community, including kafka scripts and kafka connect (default)
* `FNV-1a`, which is the standard hasher used in the Go kafka community (_e.g. Sarama producers_)

The CLI was written in Go using [spf13/cobra](https://github.com/spf13/cobra).

## Installation

It is strongly recommended that you use a released version. You can find the released binaries [here](https://github.com/ricardo-ch/kafka-topic-cloner/releases).

This project has binaries for Linux(386), and since the 0.5.0 for Windows. If your platform is not supported, you can build the project manually, or open an issue asking us to add your platform to the supported ones.

If you want to dig inside the code, you can install it with `go get`:
```sh
go get -u github.com/ricardo-ch/kafka-topic-cloner
```

## Usage

Argument        | Shorthand | Description
-----------     | --------- | -----------
brokers         | b         | Semicolon-separated list of kafka brokers `required`
from            | f         | Source topic `required`
to              | t         | Destination topic `required`
timeout         | o         | consumer timeout is ms (defaults to 10000)
default-hasher  | d         | use Sarama's FNV-1a hasher (defaults to false)
verbose         | v         | verbose mode (defaults to false)


```sh
kafka-topic-cloner --brokers my-kafka-broker-1 --from my-source-topic --to my-target-topic
```

## Disclaimer

`Kafka topic cloner` can not be considered production-ready since it does not come with any test at the moment. Use it at your own risk! :)

## Contributing

Contributions are always welcome and appreciated, the maintainers have a look at the new issues and consider every pull request.