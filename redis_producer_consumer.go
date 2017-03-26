package main

import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"sync"
	"time"
)

const TIMEOUT_SECONDS = 5

var HOST = os.Getenv("REDIS_HOST")
var PORT = os.Getenv("REDIS_PORT")
var PASSWORD = os.Getenv("REDIS_PASSWORD")

var TOTAL_BYTES_SENT = 0
var TOTAL_BYTES_READ = 0

const MIN_SEND_PERIOD = 2.0
const MAX_BUFFER_SIZE = 500

var opMutex sync.Mutex
var timerChan = make(chan int)

type RedisProducerConsumer struct {
	client      *redis.Client
	keyspace    string
	writeBuffer []string
}

func clockTick() {
	for {
		time.Sleep(MIN_SEND_PERIOD * time.Second)
		timerChan <- 1
	}
}

func NewRedisProducerConsumer(keyspace string) *RedisProducerConsumer {
	// TODO this should be a singleton
	go clockTick()
	p := &RedisProducerConsumer{
		client: redis.NewClient(
			&redis.Options{
				Addr:     fmt.Sprintf("[%s]:%s", HOST, PORT),
				Password: PASSWORD,
				DB:       0,
			},
		),
		writeBuffer: make([]string, 0),
		keyspace:    keyspace,
	}
	go p.flushOnTimerEvent()
	return p
}

func (p *RedisProducerConsumer) ClearBuffer() {
	p.client.Del(p.keyspace)
}

func (p *RedisProducerConsumer) ProduceData(data string) {
	opMutex.Lock()

	for len(p.writeBuffer) >= MAX_BUFFER_SIZE {
		p.flush()
	}
	p.writeBuffer = append(p.writeBuffer, data)
	opMutex.Unlock()

	if len(p.writeBuffer) >= MAX_BUFFER_SIZE {
		p.flush()
	}
}

func (p *RedisProducerConsumer) flushOnTimerEvent() {
	for {
		<-timerChan
		p.flush()
	}
}

func (p *RedisProducerConsumer) joinedBufferedBytes() string {
	var buffer bytes.Buffer
	for _, value := range p.writeBuffer {
		buffer.WriteString(value)
	}
	p.writeBuffer = make([]string, 0)
	return buffer.String()

}

func (p *RedisProducerConsumer) flush() {
	opMutex.Lock()
	defer opMutex.Unlock()

	if len(p.writeBuffer) == 0 {
		return
	}
	p.client.RPush(
		p.keyspace,
		p.joinedBufferedBytes(),
	)
}

func (p *RedisProducerConsumer) ConsumeData() string {
	for {
		response := p.client.BLPop(20*time.Second, p.keyspace).Val()

		if len(response) == 0 {
			continue
		}
		// key := response[0]
		ret := response[1]
		TOTAL_BYTES_READ += len(ret)
		return ret
	}
}
