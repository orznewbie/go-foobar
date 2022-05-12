package rabbitmq

import (
	"github.com/streadway/amqp"
	"strconv"
	"testing"
)

func TestMessageConsume(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		LogsDirectExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		t.Fatal(err)
	}

	normalQueue, err := ch.QueueDeclare(
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

	wrongQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := ch.QueueBind(
		normalQueue.Name,
		"info",
		LogsDirectExchange,
		false,
		nil); err != nil {
		t.Fatal(err)
	}

	if err := ch.QueueBind(
		wrongQueue.Name,
		"warning",
		LogsDirectExchange,
		false,
		nil); err != nil {
		t.Fatal(err)
	}
	if err := ch.QueueBind(
		wrongQueue.Name,
		"error",
		LogsDirectExchange,
		false,
		nil); err != nil {
		t.Fatal(err)
	}

	normalDelivery, err := ch.Consume(
		normalQueue.Name, // queue
		"",               // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              // args
	)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for dlv := range normalDelivery {
			t.Logf("queue %s receive: %s", normalQueue.Name, string(dlv.Body))
		}
	}()

	wrongDelivery, err := ch.Consume(
		wrongQueue.Name, // queue
		"",              // consumer
		true,            // auto ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for dlv := range wrongDelivery {
			t.Logf("queue %s receive: %s", wrongQueue.Name, string(dlv.Body))
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

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
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
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
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

			t.Logf("call fib(%d)..", n)
			ans := fib(n)
			t.Logf("done. fib(%d) = %d", n, ans)

			if err := ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
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
