package main

// add topic
// add subscriptions
// publish message
// subscribe
// optional test: Ack()

import "fmt"
import "./packages/broker"

func main() {

	fmt.Println("Runner started ...")

	broker.CreateTopic("topic1")

	broker.CreateTopic("topic2")

	broker.AddSubscription("topic1","sub1")
	broker.AddSubscription("topic1", "sub2")

	broker.AddSubscription("topic2","sub6")

	broker.AddSubscription("topic1", "sub3")

	broker.AddSubscription("topic2","sub7")

	t,s := broker.GetAll()

	for i:= range t {
		fmt.Println(i)
	}

	for i:= range s {
		fmt.Println(i)
	}
}