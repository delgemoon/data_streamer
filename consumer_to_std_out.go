package main

import (
	"fmt"
	//"os"
)

var toStdOutStreamer chan string = make(chan string)

func indefiniteConsumerRead(outputChannel chan string, consumer ProducerConsumer) {
	for {
		outputChannel <- consumer.ConsumeData()
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func indefiniteStdOutWrite(inputChannel chan string) {
	chunkSize := 1024
	for {
		data := []byte(<-inputChannel)
		for {
			splitPoint := min(chunkSize, len(data))
			chunk := data[:splitPoint]
			fmt.Printf("Some chunk: %d\n", len(chunk))
			data = data[splitPoint:]

			//os.Stdout.Write(chunk)

			if len(chunk) != chunkSize {
				break
			}
		}
	}
}
