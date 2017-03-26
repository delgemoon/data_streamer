package main

import (
	"encoding/base64"
	"github.com/slobdell/pubnub_communicator"
	"time"
)

const READ_BUFFER_SIZE = 200

var inputMessageChannel chan []interface{} = make(chan []interface{})

type PubnubProducerConsumer struct {
	communicator *pubnub_communicator.PubnubCommunicator
	buff         EvictingQueue
}

func NewPubnubProducerConsumer(keyspace string, read, write bool) *PubnubProducerConsumer {
	p := &PubnubProducerConsumer{
		communicator: pubnub_communicator.NewPubnubCommunicator(keyspace, read, write),
		buff:         NewEvictingQueue(READ_BUFFER_SIZE),
	}
	p.communicator.Register(inputMessageChannel)
	go p.transferReadsToBuffer()
	return p
}

func (p *PubnubProducerConsumer) transferReadsToBuffer() {
	for {
		args := <-inputMessageChannel
		message := args[0].(string)
		bytes, err := base64.StdEncoding.DecodeString(message)
		if err == nil {
			p.buff.Enqueue(string(bytes))
		}
	}
}

func (p *PubnubProducerConsumer) ClearBuffer() {
	p.buff = NewEvictingQueue(READ_BUFFER_SIZE)
}

func (p *PubnubProducerConsumer) ProduceData(data string) {
	p.communicator.SendMessage(
		base64.StdEncoding.EncodeToString([]byte(data)),
	)
}

func (p *PubnubProducerConsumer) ConsumeData() string {
	for {
		if p.buff.Length() == 0 {
			time.Sleep(100 * time.Millisecond)
		} else {
			return p.buff.Pop()
		}
	}
	panic("should not get here")
	return ""
}
