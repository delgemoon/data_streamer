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

func CreateRedisConnection(host string, port int) redis.Conn {
	networkConnection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic("Could not connect to the redis")
	}
	return redis.NewConn(
		networkConnection,
		TIMEOUT_SECONDS*time.Second,
		TIMEOUT_SECONDS*time.Second,
	)
}
