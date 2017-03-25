package main

import (
	"fmt"
	"os"
	"time"
)

var HOST = os.Getenv("REDIS_HOST")
var PORT = os.Getenv("REDIS_PORT")

type ProducerConsumer interface {
	ClearBuffer()
	ProduceData(data string)
	ConsumeData() string
}

func main() {
	conn := CreateRedisConnection(HOST, PORT)
	defer conn.Close()

	readOrWrite := os.Args[1]
	keyspace := os.Args[2]
	var producerConsumer ProducerConsumer = NewRedisProducerConsumer(conn, keyspace)

	switch readOrWrite {
	case "read":
		go indefiniteRedisRead(toStdOutStreamer, producerConsumer)
		go indefiniteStdOutWrite(toStdOutStreamer)
	case "write":
		producerConsumer.ClearBuffer()
		go indefiniteStdInRead(toRedisStreamer)
		go indefiniteRedisWrite(toRedisStreamer, producerConsumer)
	default:
		panic(fmt.Sprintf("Invalid input option: %s\n", readOrWrite))
	}

	for {
		// keep program alive without deadlock
		time.Sleep(1 * time.Second)
	}
}
