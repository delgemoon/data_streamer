package main

import (
	//"fmt"
	"os"
)

var toStdOutStreamer chan string = make(chan string)

func indefiniteRedisRead(outputChannel chan string, consumer ProducerConsumer) {
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
			data = data[splitPoint:]

			os.Stdout.Write(chunk)

			if len(chunk) != chunkSize {
				break
			}
		}
	}
}
