package main

// add topic
// add subscriptions
// publish message
// subscribe
// optional test: Ack()

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/broker"
)

func main() {

	fmt.Println("Runner started ...")

	broker.CreateTopic("topic1")

	broker.CreateTopic("topic2")

	broker.AddSubscription("topic1", "sub1")
	broker.AddSubscription("topic1", "sub2")

	broker.AddSubscription("topic2", "sub6")

	broker.AddSubscription("topic1", "sub3")

	broker.AddSubscription("topic2", "sub7")

	t, s := broker.GetAll()

	for key,topic := range t {
		fmt.Println( "map key:",key)

		fmt.Printf("for %q subscriptions are:\n", topic.TopicId)
		for index,sub := range topic.Subscriptions {
			fmt.Println(index,sub.SubscriptionId)
		}


	}
	fmt.Println("---------All subscriptions list------------")
	for key,val := range s {
		fmt.Println("key:",key,"val:",val)
	}
}
