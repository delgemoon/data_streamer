package main

import (
	"bufio"
	"io"
	"os"
)

const BUFFER_SIZE = 4 * 1024

var toRedisStreamer chan string = make(chan string)

func indefiniteStdInRead(outputChannel chan string) {
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, BUFFER_SIZE)

	for {
		bytesRead, err := reader.Read(buf[:cap(buf)])
		buf = buf[:bytesRead]
		if bytesRead == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		outputChannel <- string(buf)
	}
}

func indefiniteProducerWrite(inputChannel chan string, producer ProducerConsumer) {
	for {
		producer.ProduceData(<-inputChannel)
	}
}
