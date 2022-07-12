package rabbitmq

import (
	"strconv"
	"testing"

	"github.com/streadway/amqp"
)

const (
	TestRoutingKey = "test"

	RPCCall = "rpc_call"
)

const (
	ExchangeLogsDirect = "logs_direct"

	// The default exchange is implicitly bound to every queue, with a routing key equal to the queue name.
	// It is not possible to explicitly bind to, or unbind from the default exchange. It also cannot be deleted.
	ExchangeDefault = ""

	ExchangeDelayDirect = "exchange_delay_direct"
)

func rmqChannel() (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch, conn
}

func TestMessagePublish(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		ExchangeLogsDirect,
		amqp.ExchangeDirect,
		true,
		// If yes, the exchange will delete itself after at least
		// one queue or exchange has been bound to this one, and
		// then all queues or exchanges have been unbound.
		false,
		// If yes, clients cannot publish to this exchange directly.
		// It can only be used with exchange to exchange bindings.
		false,
		false,
		nil,
	); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		if err := ch.Publish(
			ExchangeLogsDirect,
			TestRoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("msg" + strconv.Itoa(i)),
			}); err != nil {
			t.Error(err)
		}
	}
}

// 工作队列测试
//
// Publisher将Work发布到rpc_queue中，Consumer会捞取rpc_queue中的Work，并执行相关函数
//
// 每个Work应该携带一个ReplyTo（结果队列名），用来告诉Consumer把处理后的结果投递到哪个消息队列（Publisher从该队列获取结果）
// 同时应该携带一个CorrelationId，Publisher会根据CorrelationId来过滤对应的处理结果
func TestWorkPublish(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	outputQueue, err := ch.QueueDeclare(
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

	go func() {
		for i := 1; i <= 10; i++ {
			if err := ch.Publish(
				ExchangeDefault,
				RPCCall,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: "fib",            // 消息标记，Consumer返回的消息应该携带该标记，以便Publisher识别
					ReplyTo:       outputQueue.Name, // 告诉Consumer把结果投递到哪个消息队列
					Body:          []byte(strconv.Itoa(i)),
				}); err != nil {
				t.Error(err)
			}
		}
	}()

	dlv, err := ch.Consume(
		outputQueue.Name,
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

	for d := range dlv {
		if d.CorrelationId == "fib" {
			res, err := strconv.Atoi(string(d.Body))
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		}
	}
}

func TestMessageDelayPublish(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		ExchangeDelayDirect,
		"x-delayed-message",
		true,
		// If yes, the exchange will delete itself after at least
		// one queue or exchange has been bound to this one, and
		// then all queues or exchanges have been unbound.
		false,
		// If yes, clients cannot publish to this exchange directly.
		// It can only be used with exchange to exchange bindings.
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		if err := ch.Publish(
			ExchangeDelayDirect,
			TestRoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("msg" + strconv.Itoa(i)),
				Headers: amqp.Table{
					"x-delay": (5 - i) * 1000,
				},
			}); err != nil {
			t.Error(err)
		}
	}
}
