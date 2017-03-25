package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
	"time"
)

const TIMEOUT_SECONDS = 5

type ProducerConsumer struct {
	conn     redis.Conn
	keyspace string
}

func NewProducerConsumer(conn redis.Conn, keyspace string) ProducerConsumer {
	return ProducerConsumer{
		conn:     conn,
		keyspace: keyspace,
	}
}

func (p ProducerConsumer) ClearBuffer() {
	p.conn.Do("DEL", p.keyspace)
}

func (p ProducerConsumer) ProduceData(data string) {
	p.conn.Do("RPUSH", p.keyspace, data)
}

func blpopValueToString(blpopValue interface{}) string {
	asInterface := blpopValue.([]interface{})
	return string(asInterface[1].([]byte))
}

func (p ProducerConsumer) ConsumeData() string {
	timeout := 0
	for {
		rawValue, err := p.conn.Do("BLPOP", p.keyspace, timeout)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		return blpopValueToString(rawValue)
	}
	panic("Should never reach here")
}

func CreateRedisConnection(host, port string) redis.Conn {
	url := fmt.Sprintf("%s:%s", host, port)
	networkConnection, err := net.Dial("tcp", url)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to the redis: %s", err))
	}
	return redis.NewConn(
		networkConnection,
		TIMEOUT_SECONDS*time.Second,
		TIMEOUT_SECONDS*time.Second,
	)
}
