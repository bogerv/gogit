package mqlib

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const (
	reconnectDelay     = 10 * time.Second // 连接断开后多久重连
	reconnectDetectDur = 10 * time.Second
)

var (
	errAlreadyClosed = errors.New("already closed: not connected to the AMQP server")
)

var _connection *amqp.Connection

type Amqp struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	done       chan bool
	changeConn chan struct{}
	connNotify chan *amqp.Error
	chanNotify chan *amqp.Error

	isConnected bool
	hasConsumer bool
}

func Init() (err error) {
	// amqp://user:pwd@host:port/
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		"", "",
		"", "",
	)
	if _connection, err = amqp.Dial(uri); err != nil {
		return
	}

	return nil
}

func New() (*Amqp, error) {
	entity := &Amqp{
		conn:       _connection,
		done:       make(chan bool),
		changeConn: make(chan struct{}, 1),
	}
	go entity.doReconnect()

	channel, err := _connection.Channel()
	if err != nil {
		_ = _connection.Close()
		return nil, err
	}

	entity.channel = channel
	entity.isConnected = true

	return entity, nil
}

func (slf *Amqp) doReconnect() {
	for {
		if !slf.isConnected {
			log.Println("Attempting to connect")
			var (
				connected = false
				err       error
			)

			for cnt := 0; !connected; cnt++ {
				if connected, err = slf.connect(); err != nil {
					log.Printf("Failed to connect: %s.\n", err)
				}
				if !connected {
					log.Printf("Retrying... %d\n", cnt)
				}
				time.Sleep(reconnectDelay)
			}
		}

		select {
		case <-slf.done:
			println("evt `slf.done` triggered")
			return
		case err := <-slf.chanNotify:
			log.Printf("channel close notify: %v", err)
			slf.isConnected = false
		case err := <-slf.connNotify:
			log.Printf("conn close notify: %v", err)
			slf.isConnected = false
		}
		time.Sleep(reconnectDetectDur)
	}
}

func (slf *Amqp) connect() (bool, error) {
	err := Init()
	if err != nil {
		return false, err
	}

	ch, err := _connection.Channel()
	if err != nil {
		return false, err
	}

	slf.isConnected = true
	slf.changeConnection(_connection, ch)
	return true, nil
}

// 监听Rabbit channel的状态
func (slf *Amqp) changeConnection(conn *amqp.Connection, channel *amqp.Channel) {
	slf.conn = conn
	slf.connNotify = make(chan *amqp.Error, 1)
	slf.conn.NotifyClose(slf.connNotify)

	slf.channel = channel
	slf.chanNotify = make(chan *amqp.Error, 1)
	slf.channel.NotifyClose(slf.chanNotify)

	// TOFIX: only producer will be blocked here
	if slf.hasConsumer {
		// true: cause only consumer will be  notify for now.
		slf.changeConn <- struct{}{}
	}
}

func (slf *Amqp) Close() error {
	if !slf.isConnected {
		return errAlreadyClosed
	}

	err := slf.channel.Close()
	if err != nil {
		return err
	}
	err = slf.conn.Close()
	if err != nil {
		return err
	}

	close(slf.done)
	slf.isConnected = false
	return nil
}

func (slf *Amqp) Send(exchange, key string, data []byte, mandatory, immediate bool) error {
	return slf.channel.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
}
