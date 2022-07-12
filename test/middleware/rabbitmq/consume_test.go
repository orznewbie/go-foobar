package rabbitmq

import (
	"strconv"
	"testing"

	"github.com/streadway/amqp"
)

func TestMessageConsume(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		ExchangeLogsDirect,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		t.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		// If "", an auto-generate name will be used
		"",
		false,
		// delete when unused
		false,
		// Exclusive queues are only accessible by the connection that declares them and
		// will be deleted when the connection closes
		true,
		// When noWait is true, the queue will assume to be declared on the server.  A
		// channel exception will arrive if the conditions are met for existing queues
		// or attempting to modify an existing queue from a different connection.
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.QueueBind(
		queue.Name,
		TestRoutingKey,
		ExchangeLogsDirect,
		false,
		nil); err != nil {
		t.Fatal(err)
	}

	delivery, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for dlv := range delivery {
			t.Logf("queue(name=%s) receive: %s", queue.Name, string(dlv.Body))
		}
	}()

	var forever chan struct{}
	<-forever
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func TestWorkConsume(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	inputQueue, err := ch.QueueDeclare(
		RPCCall,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := ch.Qos(
		1,
		0,
		false); err != nil {
		t.Fatal(err)
	}

	dlv, err := ch.Consume(
		inputQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		// TODO 为什么autoAck=false会阻塞在for..range?
		for d := range dlv {
			n, err := strconv.Atoi(string(d.Body))
			if err != nil {
				t.Error(err)
				return
			}

			ans := fib(n)
			if err := ch.Publish(
				ExchangeDefault,
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(ans)),
				}); err != nil {
				t.Fatal(err)
			}

			d.Ack(false)
		}
	}()

	var forever chan struct{}
	<-forever
}

func TestMessageDelayConsume(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		ExchangeDelayDirect,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	); err != nil {
		t.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		// If "", an auto-generate name will be used
		"",
		false,
		// delete when unused
		false,
		// Exclusive queues are only accessible by the connection that declares them and
		// will be deleted when the connection closes
		true,
		// When noWait is true, the queue will assume to be declared on the server.  A
		// channel exception will arrive if the conditions are met for existing queues
		// or attempting to modify an existing queue from a different connection.
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.QueueBind(
		queue.Name,
		TestRoutingKey,
		ExchangeDelayDirect,
		false,
		nil); err != nil {
		t.Fatal(err)
	}

	delivery, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for dlv := range delivery {
			t.Logf("queue(name=%s) receive: %s", queue.Name, string(dlv.Body))
		}
	}()

	var forever chan struct{}
	<-forever
}
