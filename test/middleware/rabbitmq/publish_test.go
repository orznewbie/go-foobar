package rabbitmq

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

const (
	LogsDirectExchange = "logs_direct"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var logLevel = []string{"info", "warning", "error"}

func randomLogLevel() string {
	return logLevel[rand.Intn(len(logLevel))]
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

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

// 消息投递测试
// 1. Fanout 2. Direct 3. Topic 4. Headers
func TestMessagePublish(t *testing.T) {
	ch, conn := rmqChannel()
	defer func() {
		conn.Close()
		ch.Close()
	}()

	if err := ch.ExchangeDeclare(
		LogsDirectExchange,
		"direct",
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

	for i := 0; i < 20; i++ {
		level := randomLogLevel()
		if err := ch.Publish(
			LogsDirectExchange, // exchange
			level,              // routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(level + " log"),
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		t.Fatal(err)
	}

	corrId := randomString(32)

	go func() {
		for i := 1; i <= 10; i++ {
			if err := ch.Publish(
				"",          // exchange
				"rpc_queue", // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: corrId, // 消息标记，Consumer返回的消息应该携带该标记，以便Publisher识别
					ReplyTo:       q.Name, // 告诉Consumer把结果投递到哪个消息队列
					Body:          []byte(strconv.Itoa(i)),
				}); err != nil {
				t.Error(err)
			}
		}
	}()

	dlv, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		t.Fatal(err)
	}

	for d := range dlv {
		if corrId == d.CorrelationId {
			res, err := strconv.Atoi(string(d.Body))
			if err != nil {
				t.Fatal(err)
			}
			t.Log(res)
		}
	}
}
