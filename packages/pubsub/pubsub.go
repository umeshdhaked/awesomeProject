package pubsub

import (
	"errors"
	"fmt"
	"log"
	"time"
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
	messageIdTracker   int
	flag               bool
	ch                 chan *Message
	topics             topics
	subscriptionTopics subscriptionTopics
}

const retryTime = time.Second * 15


func NewPubSub() IPubSub {
	var pubsub *PubSub = &PubSub{0, true, make(chan *Message, 50), topics{topicsMap: make(map[string]*topic)}, subscriptionTopics{subscriptionTopicMap: make(map[string]string)}}
	return pubsub
}


func (p *PubSub) Subscribe(subscriptionID string, subscriberFunc func(msg Message)) (bool, error) {

	p.subscriptionTopics.subscriptionTopicMapMutex.RLock()
	topicId, ok := p.subscriptionTopics.subscriptionTopicMap[subscriptionID]
	p.subscriptionTopics.subscriptionTopicMapMutex.RUnlock()

	if ok {
		p.topics.topicMutex.RLock()
		var topicVar *topic = p.topics.topicsMap[topicId]
		p.topics.topicMutex.RUnlock()

		topicVar.subscriptionsMutex.RLock()
		var subscriptionVar *subscription = topicVar.subscriptions[subscriptionID]
		topicVar.subscriptionsMutex.RUnlock()


		subscriberVar := subscriber(subscriberFunc)
		subscriptionVar.addSubscriber(&subscriberVar)
		return true, nil
	} else {
		log.Println("Subscribe()-> Subscription doesn't exist, subscriptionID:",subscriptionID)

		return true, errors.New("this subscription doesn't exist")
	}

}

func (p *PubSub) UnSubscribe(subscriptionID string) (bool, error) {

	p.subscriptionTopics.subscriptionTopicMapMutex.RLock()
	topicId, ok := p.subscriptionTopics.subscriptionTopicMap[subscriptionID]
	p.subscriptionTopics.subscriptionTopicMapMutex.RUnlock()


	if ok {
		p.topics.topicMutex.RLock()
		var topicVar *topic = p.topics.topicsMap[topicId]
		p.topics.topicMutex.RUnlock()

		topicVar.subscriptionsMutex.RLock()
		var subscriptionVar *subscription = topicVar.subscriptions[subscriptionID]
		topicVar.subscriptionsMutex.RUnlock()

		subscriptionVar.removeSubscriber()

		return true, nil
	} else {
		log.Println("UnSubscribe()-> Subscription doesn't exist, subscriptionID:",subscriptionID)
		return true, errors.New("this subscription doesn't exist")
	}
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

		(*p).topics.topicMutex.RLock()
		var topicVar *topic = (*p).topics.topicsMap[msg.topicId]
		(*p).topics.topicMutex.RUnlock()

		var subsObjs map[string]*subscription = topicVar.subscriptions

		topicVar.subscriptionsMutex.RLock()
		for _, val := range subsObjs {
			go val.sendMessage(msg)
		}
		topicVar.subscriptionsMutex.RUnlock()
	}

}

func (p *PubSub) Ack(msgId int, subId string) (bool, error) {

	p.subscriptionTopics.subscriptionTopicMapMutex.RLock()
	topicId, ok := p.subscriptionTopics.subscriptionTopicMap[subId]
	p.subscriptionTopics.subscriptionTopicMapMutex.RUnlock()

	if ok {

		p.topics.topicMutex.RLock()
		var topicVar *topic = p.topics.topicsMap[topicId]
		p.topics.topicMutex.RUnlock()

		topicVar.subscriptionsMutex.RLock()
		var subscriptionVar *subscription = topicVar.subscriptions[subId]
		topicVar.subscriptionsMutex.RUnlock()


		subscriptionVar.pendingMapMutex.Lock()
		_, ok1 := subscriptionVar.pendingAckMsg[msgId]
		if ok1 {
			delete(subscriptionVar.pendingAckMsg, msgId)
			subscriptionVar.pendingMapMutex.Unlock()
			fmt.Printf("Ack()-> Message id %v has been Acknowledge by %v \n", msgId, subId)
			return true, nil
		}
		subscriptionVar.pendingMapMutex.Unlock()
		fmt.Printf("Ack()-> Message is already Acknowledged or wrong messageId: %v in subscription: %v \n", msgId, subId)
		return false, errors.New("message is already Ack or wrong subscriptionId")
	} else {
		fmt.Printf("Ack()-> subscription %v doesn't exist", subId)
		return false,errors.New("subscription doesn't exist")
	}

}
