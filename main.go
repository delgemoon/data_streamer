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
	// TODO do this better in the redis constructor
	//conn := CreateRedisConnection(HOST, PORT)
	//defer conn.Close()

	readOrWrite := os.Args[1]
	keyspace := os.Args[2]
	// var producerConsumer ProducerConsumer = NewRedisProducerConsumer(conn, keyspace)
	read := false
	write := false
	if readOrWrite == "read" {
		read = true
	} else {
		write = true
	}
	var producerConsumer ProducerConsumer = NewPubnubProducerConsumer(keyspace, read, write)

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
