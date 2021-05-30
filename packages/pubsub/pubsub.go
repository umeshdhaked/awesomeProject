package pubsub

import (
	"errors"
	"log"
	"sync"
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
const retryMaxCount = 5

func NewPubSub() IPubSub {
	var pubsub *PubSub = &PubSub{0, true, make(chan *Message, 50), topics{topicsMap: new(sync.Map)}, subscriptionTopics{subscriptionTopicMap: new(sync.Map)}}
	return pubsub
}

func (p *PubSub) Subscribe(subscriptionID string, subscriberFunc func(msg Message)) (bool, error) {

	topicId, ok := p.subscriptionTopics.subscriptionTopicMap.Load(subscriptionID)

	if ok {
		topicVar, _ := p.topics.topicsMap.Load(topicId)

		subscriptionVar, _ := topicVar.(*topic).subscriptions.Load(subscriptionID)

		subscriberVar := subscriber(subscriberFunc)
		subscriptionVar.(*subscription).addSubscriber(&subscriberVar)
		return true, nil
	} else {
		log.Println("Subscribe()-> Subscription doesn't exist, subscriptionID:", subscriptionID)

		return true, errors.New("this subscription doesn't exist")
	}

}

func (p *PubSub) UnSubscribe(subscriptionID string) (bool, error) {

	topicId, ok := p.subscriptionTopics.subscriptionTopicMap.Load(subscriptionID)

	if ok {
		topicVar, _ := p.topics.topicsMap.Load(topicId.(string))

		subscriptionVar, _ := topicVar.(*topic).subscriptions.Load(subscriptionID)

		subscriptionVar.(*subscription).removeSubscriber()

		return true, nil
	} else {
		log.Println("UnSubscribe()-> Subscription doesn't exist, subscriptionID:", subscriptionID)
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

		topicVar, ok := (*p).topics.topicsMap.Load(msg.topicId)

		if ok {
			var size = 0
			topicVar.(*topic).subscriptions.Range(func(_, value interface{}) bool {
				go value.(*subscription).sendMessage(msg)
				size++
				return true
			})

			if size == 0 {
				log.Printf("pushMessage()-> this topic has no subscription")
			}
		} else {
			log.Printf("pushMessage()-> this topic doesn't exists, unknown topic:" + msg.topicId)
		}
	}

}

func (p *PubSub) Ack(msgId int, subId string) (bool, error) {

	topicId, ok := p.subscriptionTopics.subscriptionTopicMap.Load(subId)

	if ok {

		topicVar, _ := p.topics.topicsMap.Load(topicId.(string))

		subscriptionVar, _ := topicVar.(*topic).subscriptions.Load(subId)

		_, ok1 := subscriptionVar.(*subscription).pendingAckMsg.Load(msgId)
		if ok1 {
			subscriptionVar.(*subscription).pendingAckMsg.Delete(msgId)
			log.Printf("Ack()-> Message id %v has been Acknowledge by %v \n", msgId, subId)
			return true, nil
		}
		log.Printf("Ack()-> Message is already Acknowledged or wrong messageId: %v in subscription: %v \n", msgId, subId)
		return false, errors.New("message is already Ack or wrong subscriptionId")
	} else {
		log.Printf("Ack()-> subscription %v doesn't exist", subId)
		return false, errors.New("subscription doesn't exist")
	}

}
