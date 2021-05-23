package main

// add topic
// add subscriptions
// publish message
// subscribe
// optional test: Ack()

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/broker"
	"time"
)

func createTpc(id string) {
	broker.CreateTopic(id)
}

func createSub(id string) {
	broker.AddSubscription("topic1",id)
}

func main() {

	fmt.Println("Runner started ...")

	broker.CreateTopic("topic1")

	for i := 0 ; i<100 ; i++ {
		go createTpc(fmt.Sprintf("%v", i))
	}

	for i := 0 ; i<100 ; i++ {
		go createSub(fmt.Sprintf("%v", i))
	}

	time.Sleep(2*time.Second)

	broker.AddSubscription("topic1", "sub1")


	broker.Subscribe("sub1",messageReceiverFunc)

	broker.Publish("topic1","Hello Subscriber1, Are ya winning son?")
	broker.Publish("topic1","Hello Subscriber2, Are ya winning son?")
	broker.Publish("topic1","Hello Subscriber3, Are ya winning son?")
	broker.Publish("topic1","Hello Subscriber4, Are ya winning son?")
	broker.Publish("topic1","Hello Subscriber5, Are ya winning son?")
	broker.Publish("topic1","Hello Subscriber6, Are ya winning son?")


}

var messageReceiverFunc = func(message broker.Message) {
	fmt.Println(message.MessageId,message.Data)
}