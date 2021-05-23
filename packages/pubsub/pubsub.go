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

	Ack(msgId int, subId string)
}

type PubSub struct {
	messageIdTracker int
	flag             bool
	ch               chan *Message
	topics           Topics
	subscriptions    Subscriptions
}

var pubsub = &PubSub{0, true, make(chan *Message, 10), Topics{topicsMap: make(map[string]*Topic)}, Subscriptions{subscriptionMap: make(map[string]*Subscription)}}

func GetPubSub() IPubSub {
	return pubsub
}

func GetAllTopicsAndSubscriptions() (map[string]*Topic, map[string]*Subscription) {
	return pubsub.topics.topicsMap, pubsub.subscriptions.subscriptionMap
}

func (p *PubSub) Subscribe(subscriptionID string, subscriberFunc func(msg Message)) {

	var subscription = pubsub.subscriptions.subscriptionMap[subscriptionID]

	subscriber := Subscriber(subscriberFunc)

	subscription.addSubscriber(&subscriber)

}

func (p *PubSub) UnSubscribe(subId string) {
	var subscription = pubsub.subscriptions.subscriptionMap[subId]
	subscription.removeSubscriber()
}

func (p *PubSub) Publish(topicId, message string) {

	p.messageIdTracker++
	messageObj := &Message{MessageId: p.messageIdTracker, TopicId: topicId, Data: message}

	if p.flag {
		go pushMessage(p.ch)
		p.flag = false
	}
	p.ch <- messageObj
}

func pushMessage(ch chan *Message) {

	for {
		msg := <-ch
		var topic *Topic = pubsub.topics.topicsMap[msg.TopicId]

		var subsObjs []*Subscription = topic.Subscriptions

		for _, val := range subsObjs {
			go val.sendMessage(msg)
		}
	}

}

func (p *PubSub) Ack(msgId int, subId string) {

	p.subscriptions.subMutex.RLock()
	subs, ok := p.subscriptions.subscriptionMap[subId]
	p.subscriptions.subMutex.RUnlock()

	if ok {
		subs.pendingMapMutex.Lock()
		_, ok1 := subs.pendingAck[msgId]
		if ok1 {
			delete(subs.pendingAck, msgId)
			fmt.Printf("Message id %q has been Acknowledge by %q \n", msgId, subId)
		} else {
			fmt.Printf("Message is alrady Ack or wrong messageId: %q in subscription: %q \n", msgId, subId)
		}
		subs.pendingMapMutex.Unlock()
	} else {
		fmt.Printf("Subscription %q doesn't exist", subId)
	}

}
