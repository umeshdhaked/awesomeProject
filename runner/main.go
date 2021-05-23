package main

// add topic
// add subscriptions
// publish message
// subscribe
// optional test: Ack()

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"github.com/umeshdhaked/awesomeProject/packages/subscribers"
	"github.com/umeshdhaked/awesomeProject/runner/api"
)

var pubsubObj = pubsub.GetPubSub()

func createTpc(id string) {
	pubsubObj.CreateTopic(id)
}

func createSub(id string) {
	pubsubObj.AddSubscription("topic1",id)
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

	pubsubObj.CreateTopic("topic1")
	pubsubObj.AddSubscription("topic1", "sub1")
	pubsubObj.AddSubscription("topic1", "sub2")

	pubsubObj.Subscribe("sub1",subscribers.SubscriberTypeA)
	pubsubObj.Subscribe("sub2",subscribers.SubscriberTypeB)

	pubsubObj.Publish("topic1","Hello Subscriber1, Are ya winning son?")
	pubsubObj.Publish("topic1","Hello Subscriber2, Are ya winning son?")
	pubsubObj.Publish("topic1","Hello Subscriber3, Are ya winning son?")

	//pubsubObj.Unsubscribe("sub1")
	pubsubObj.Publish("topic1","Hello Subscriber4, Are ya winning son?")
	pubsubObj.Publish("topic1","Hello Subscriber5, Are ya winning son?")
	pubsubObj.Publish("topic1","Hello Subscriber6, Are ya winning son?")


	api.HandleRequest()

}