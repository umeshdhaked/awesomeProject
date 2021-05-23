package main

// add topic
// add subscriptions
// publish message
// subscribe
// optional test: Ack()

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"github.com/umeshdhaked/awesomeProject/runner/api"
	"log"
)

func createTpc(id string) {
	pubsub.CreateTopic(id)
}

func createSub(id string) {
	pubsub.AddSubscription("topic1",id)
}

func main() {

	fmt.Println("Runner started ...")

	//for i := 0 ; i<100 ; i++ {
	//	go createTpc(fmt.Sprintf("%v", i))
	//}
	//
	//for i := 0 ; i<100 ; i++ {
	//	go createSub(fmt.Sprintf("%v", i))
	//}
	//time.Sleep(2*time.Second)

	pubsub.CreateTopic("topic1")
	pubsub.AddSubscription("topic1", "sub1")
	pubsub.Subscribe("sub1",messageReceiverFunc1)

	pubsub.Publish("topic1","Hello Subscriber1, Are ya winning son?")
	pubsub.Publish("topic1","Hello Subscriber2, Are ya winning son?")
	pubsub.Publish("topic1","Hello Subscriber3, Are ya winning son?")

	//pubsub.Unsubscribe("sub1")
	pubsub.Publish("topic1","Hello Subscriber4, Are ya winning son?")
	pubsub.Publish("topic1","Hello Subscriber5, Are ya winning son?")
	pubsub.Publish("topic1","Hello Subscriber6, Are ya winning son?")


	api.HandleRequest()

}

var messageReceiverFunc1 = func(message pubsub.Message) {
	fmt.Println("messageId:",message.MessageId,"| message:",message.Data)
}

var messageReceiverFunc2 = func(message pubsub.Message) {
	log.Println("messageId:",message.MessageId,"| message:",message.Data)
}