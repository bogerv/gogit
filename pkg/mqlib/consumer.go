package mqlib

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

// ConsumeHandler .
type ConsumeHandler func(d amqp.Delivery) error

func (slf *Amqp) Consume(
	queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table,
	handler ConsumeHandler,
) {
	var (
		delivery <-chan amqp.Delivery
		err      error
	)

	slf.hasConsumer = true

	for {
		select {
		case <-slf.changeConn:
			log.Println("evt 'changeConn' triggered.")
			if delivery, err = slf.channel.Consume(
				queue,     // queue
				consumer,  // consumer
				autoAck,   // auto-ack
				exclusive, // exclusive
				noLocal,   // no-local
				noWait,    // no-wait
				args,      // args
			); err != nil {
				log.Println("could not start consuming with error: ", err)
				break
			}
			log.Println("initial comsumer finished")
		default:
			if !slf.isConnected || delivery == nil {
				// true: wrapper has not connected or consumer has not initialized
				// must to wait `changeConn` evt
				time.Sleep(1 * time.Second)
				break
			}
			// delivery will be closed, then this `range` will be finished
			for d := range delivery {
				if err := handler(d); err != nil {
					log.Printf("could not consume message: %v with error: %v", d, err)
				}
			}
		}
	}
}
