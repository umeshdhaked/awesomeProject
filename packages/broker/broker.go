package broker

import "fmt"

type PubSub interface {
	CreateTopic()

	AddSubscription()

	Publish()

	Subscribe()

	Ack()
}


//  Create subscription
var topics map[string]*Topic = make(map[string]*Topic)

func CreateTopic(topicName string) bool {
	val, ok := topics[topicName]

	if ok {
		fmt.Println("Topic Already Exists", val)
		return false
	} else {
		topics[topicName] = &Topic{TopicId: topicName}
		return true
	}
}

//  Add subscription
var subscriptions map[string]*Subscription = make(map[string]*Subscription)

func AddSubscription(topicID string, subName string) bool {
	val, ok := subscriptions[subName]

	if ok {
		fmt.Println("Subscription Already Exists", val)
		return false
	} else {
		topic, ok := topics[topicID]
		if !ok {
			fmt.Println("Topic doesn't Exists")
			return false
		}

		subscriptions[subName] = &Subscription{SubscriptionId: subName}
		topic.Subscriptions = append(topic.Subscriptions, subscriptions[subName] )
		return true
	}
}


func GetAll() (map[string]*Topic, map[string]*Subscription){
	return topics,subscriptions
}

