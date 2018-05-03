# go-kafka-connect
Go Cobra CLI that clones events between two topics of a unique kafka cluster

## Clone

### Arguments

Argument    | Shorthand | Description
----------- | --------- | -----------
broker-list | b         | Semicolon-separated list of kafka brokers
from        | f         | Original topic
to          | t         | Destination topic

```bash
kafka-topic-cloner clone --bootstrap-server my-kafka-broker-1;my-kafka-broker-2 --from my-source-topic --to my-target-topic
```