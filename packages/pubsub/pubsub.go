package pubsub

import (
	"errors"
	"fmt"
)

type IPubSub interface {
	CreateTopic(topicName string) (bool, error)

	DeleteTopic(TopicID string) (bool, error)

	AddSubscription(topicID, subName string) (bool, error)

	DeleteSubscription(SubscriptionID string) (bool, error)

	Subscribe(subscriptionID string, subscriberFunc func(msg Message)) (bool, error)

	UnSubscribe(subId string) (bool, error)

	Publish(topicId, message string) (bool, error)

	Ack(msgId int, subId string) (bool, error)
}

type PubSub struct {
	messageIdTracker int
	flag             bool
	ch               chan *Message
	topics           topics
	subscriptions    subscriptions
}

var pubsub *PubSub = &PubSub{0, true, make(chan *Message, 50), topics{topicsMap: make(map[string]*topic)}, subscriptions{subscriptionMap: make(map[string]*subscription)}}

func GetPubSub() IPubSub {
	return pubsub
}


func (p *PubSub) Subscribe(subscriptionID string, subscriberFunc func(msg Message)) (bool, error) {

	var subscription = p.subscriptions.subscriptionMap[subscriptionID]

	subscriber := subscriber(subscriberFunc)

	subscription.addSubscriber(&subscriber)

	return true, nil

}

func (p *PubSub) UnSubscribe(subId string) (bool, error) {
	var subscription = p.subscriptions.subscriptionMap[subId]
	subscription.removeSubscriber()

	return true, nil
}

func (p *PubSub) Publish(topicId, message string) (bool, error) {

	p.messageIdTracker++
	messageObj := &Message{messageId: p.messageIdTracker, topicId: topicId, data: message}

	if p.flag {
		go pushMessage(p)
		p.flag = false
	}
	p.ch <- messageObj

	return true, nil
}

func pushMessage(p *PubSub) {

	for {
		msg := <-(*p).ch
		var topicVar *topic = (*p).topics.topicsMap[msg.topicId]

		var subsObjs []*subscription = topicVar.subscriptions

		for _, val := range subsObjs {
			go val.sendMessage(msg)
		}
	}

}

func (p *PubSub) Ack(msgId int, subId string) (bool, error) {

	p.subscriptions.subMutex.RLock()
	subs, ok := p.subscriptions.subscriptionMap[subId]
	p.subscriptions.subMutex.RUnlock()

	if ok {
		subs.pendingMapMutex.Lock()
		_, ok1 := subs.pendingAckMsg[msgId]
		if ok1 {
			delete(subs.pendingAckMsg, msgId)
			subs.pendingMapMutex.Unlock()
			fmt.Printf("Ack()-> Message id %v has been Acknowledge by %v \n", msgId, subId)
			return true, nil
		}
		subs.pendingMapMutex.Unlock()
		fmt.Printf("Ack()-> Message is already Acknowledged or wrong messageId: %v in subscription: %v \n", msgId, subId)
		return false, errors.New("message is already Ack or wrong subscriptionId")
	} else {
		fmt.Printf("Ack()-> subscription %v doesn't exist", subId)
		return false,errors.New("subscription doesn't exist")
	}

}
