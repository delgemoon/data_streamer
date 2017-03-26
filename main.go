package main

import (
	"fmt"
	"os"
	"time"
)

type ProducerConsumer interface {
	ClearBuffer()
	ProduceData(data string)
	ConsumeData() string
}

func main() {
	readOrWrite := os.Args[1]
	keyspace := os.Args[2]
	var producerConsumer ProducerConsumer = NewRedisProducerConsumer(keyspace)
	//var producerConsumer ProducerConsumer = NewPubnubProducerConsumer(keyspace, read, write)

	switch readOrWrite {
	case "read":
		go indefiniteConsumerRead(toStdOutStreamer, producerConsumer)
		go indefiniteStdOutWrite(toStdOutStreamer)
	case "write":
		producerConsumer.ClearBuffer()
		go indefiniteStdInRead(toRedisStreamer)
		go indefiniteProducerWrite(toRedisStreamer, producerConsumer)
	default:
		panic(fmt.Sprintf("Invalid input option: %s\n", readOrWrite))
	}

	for {
		// keep program alive without deadlock
		time.Sleep(1 * time.Second)
	}
}
