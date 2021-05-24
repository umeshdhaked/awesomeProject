package main

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"math/rand"
	"time"
)

var pubSubObj pubsub.IPubSub = pubsub.NewPubSub()
var topic1, topic2 string = "topic1", "topic2"
var topic1_sub1, topic1_sub2, topic2_sub1 string = "topic1_sub1", "topic1_sub2", "topic2_sub1"

func createTopic(topicId string) {
	ok, err := pubSubObj.CreateTopic(topicId)
	if ok {
		fmt.Println("Created topic {topicID : "+topicId+"}")
	}else {
		fmt.Println(err)
	}
}

func addSubscription(topicId string, subscriptionId string)  {
	ok, err := pubSubObj.AddSubscription(topicId, subscriptionId)
	if ok {
		fmt.Println("Created subscription {topicID: "+topicId+", subscriptionID: "+subscriptionId+"}")
	}else {
		fmt.Println(err)
	}
}

func subscribe(subscriptionId string, subsFunc func(msg pubsub.Message)) {
	ok, err := pubSubObj.Subscribe(subscriptionId, subsFunc)
	if ok {
		fmt.Println("Subscribed {SubscriptionID: "+subscriptionId+"}")
	}else {
		fmt.Println(err)
	}
}

func publish(topicID, msg string) {
	ok, err := pubSubObj.Publish(topicID, msg)
	if ok {
		fmt.Println("Published {topicID: "+topicID+", Message: "+msg + "}")
	}else {
		fmt.Println(err)
	}
}

func main() {
	fmt.Printf("This is the simulation of library with default hardcoded configurations :) \n\n\n\n ")

	createTopic(topic1)
	createTopic(topic2)

	addSubscription(topic1, topic1_sub1)
	addSubscription(topic1, topic1_sub2)
	addSubscription(topic1, topic1_sub1) // duplication subscription, should get error
	addSubscription(topic2, topic2_sub1)

	subscribe(topic1_sub1, subscriberFuncA)
	subscribe(topic1_sub2, subscriberFuncB)
	subscribe(topic2_sub1, subscriberFuncC)

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		publish(topic1, fmt.Sprintf("random Message: %v", seededRand.Int()))
	}

	for i := 0; i < 5; i++ {
		publish(topic2, fmt.Sprintf("random Message: %v", seededRand.Int()))
	}

	time.Sleep(time.Second * 20)

	pubSubObj.DeleteTopic(topic1)
	time.Sleep(time.Second * 3)
	pubSubObj.Publish(topic1, fmt.Sprintf("random Message: %v", seededRand.Int()))

	time.Sleep(time.Minute * 2)

	fmt.Println("Exiting after completion ... ")
}

// this function will ACK after receiving msg
func subscriberFuncA(msg pubsub.Message) {
	defer pubSubObj.Ack(msg.MessageId(), topic1_sub1)
	fmt.Println("SubscriberTypeA, Received message : ", msg.Data())
}

// this function will ACK after 18 seconds of receiving msg
func subscriberFuncB(msg pubsub.Message) {
	defer pubSubObj.Ack(msg.MessageId(), topic1_sub2)
	fmt.Println("SubscriberTypeB, Received message :", msg.Data())
	time.Sleep(time.Second * 18)
}

// this function will ACK after receiving msg
func subscriberFuncC(msg pubsub.Message) {
	defer pubSubObj.Ack(msg.MessageId(), topic2_sub1)
	fmt.Println("SubscriberTypeC, Received message :", msg.Data())
}
