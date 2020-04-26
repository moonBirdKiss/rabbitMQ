package main

import (
	"github.com/kataras/iris/v12"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	app := iris.New()

	// amqp info
	amqp_url := "amqp://admin:admin@192.168.66.90:5672/"


	log.Println("1.0 server start to work")

	app.Get("/test", func(ctx iris.Context) {
		ctx.WriteString("server work properly")
	})

	// sender
	app.Post("/send", func(ctx iris.Context){
		queueName := ctx.PostValue("name")
		data := ctx.PostValue("data")


		// send msg to rabbitMQ
		conn, err := amqp.Dial(amqp_url)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			queueName, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(data),
			})

		log.Println("success send msg")
		ctx.WriteString(queueName + data)
	})


	// save data
	savedQueue := make(map[string] chan []byte)

	// receiver
	app.Post("/receive", func(ctx iris.Context) {
		// get value
		queueName := ctx.PostValue("name")

		_, exist := savedQueue[queueName]

		// if the queue needs init,
		if exist == false{

			log.Println("2.0 start init queue in rabbit_MQ Receiver")

			queueData := make(chan []byte, 10)
			savedQueue[queueName] = queueData

			// read data from queue
			// and this func will exe forever
			go func() {
				conn, err := amqp.Dial(amqp_url)
				failOnError(err, "Failed to connect to RabbitMQ")
				defer conn.Close()

				ch, err := conn.Channel()
				failOnError(err, "Failed to open a channel")
				defer ch.Close()

				q, err := ch.QueueDeclare(
					queueName, // name
					false,   // durable
					false,   // delete when unused
					false,   // exclusive
					false,   // no-wait
					nil,     // arguments
				)
				failOnError(err, "Failed to declare a queue")
				// it return a channel
				// so we can do a lot of things
				msgs, err := ch.Consume(
					q.Name, // queue
					"",     // consumer
					true,   // auto-ack
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
				)
				failOnError(err, "Failed to register a consumer")

				forever := make(chan bool)

				go func() {
					for d := range msgs {
						log.Printf("Received a message: %s", d.Body)
						savedQueue[queueName] <- d.Body
					}
				}()
				log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

				// and this for infinet loop
				<-forever
			}()
		}

		// get data from queue
		data := <- savedQueue[queueName]
		ctx.Write(data)
	})

	app.Run(iris.Addr("0.0.0.0:8090"), iris.WithoutServerError(iris.ErrServerClosed))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}