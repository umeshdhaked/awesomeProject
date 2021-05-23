package pubsub

import "fmt"

type IPubSub interface {
	CreateTopic(topicName string) bool

	DeleteTopic(TopicID string)

	AddSubscription(topicID, subName string) bool

	DeleteSubscription(SubscriptionID string)

	Subscribe(subscriptionID string, subscriberFunc func(msg Message))

	UnSubscribe(subId string)

	Publish(topicId, message string)

	Ack(msgId, subId string)
}

type PubSub struct {
	messageIdTracker int
	flag             bool
	ch               chan Message
}

var pubsub *PubSub = &PubSub{0, true, make(chan Message, 10)}

func GetPubSub() IPubSub{
	return pubsub
}

func GetAllTopicsAndSubscriptions() (map[string]*Topic, map[string]*Subscription) {
	return topics.topicsMap, subscriptions.subscriptionMap
}

func (p *PubSub) Subscribe(subscriptionID string, subscriberFunc func(msg Message)) {

	var subscription = subscriptions.subscriptionMap[subscriptionID]

	subscriber := Subscriber(subscriberFunc)

	subscription.addSubscriber(&subscriber)

}

func (p *PubSub) UnSubscribe(subId string) {
	var subscription = subscriptions.subscriptionMap[subId]
	subscription.removeSubscriber()
}

func (p *PubSub) Publish(topicId, message string) {

	p.messageIdTracker++
	messageObj := Message{p.messageIdTracker, topicId, message}

	if p.flag {
		go pushMessage(p.ch)
		p.flag = false
	}
	p.ch <- messageObj
}

func pushMessage(ch chan Message) {

	for {
		msg := <-ch
		var topic = topics.topicsMap[msg.TopicId]

		var subsObjs = topic.Subscriptions

		for _, val := range subsObjs {
			go val.sendMessage(&msg)
		}
	}

}

func (p *PubSub) Ack(msgId, subId string) {
	fmt.Println("Not yet implemented")
}