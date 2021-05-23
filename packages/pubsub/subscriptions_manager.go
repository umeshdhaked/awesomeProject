package pubsub

import (
	"fmt"
	"log"
	"sync"
)

type Subscriptions struct {
	subscriptionMap map[string]*Subscription
	subMutex        sync.RWMutex
}

//  Add subscription

func (p *PubSub) AddSubscription(topicID, subName string) bool {
	pubsub.subscriptions.subMutex.Lock()
	defer pubsub.subscriptions.subMutex.Unlock()

	val, ok := pubsub.subscriptions.subscriptionMap[subName]
	if ok {
		fmt.Println("Subscription Already Exists", val)
		return false
	} else {
		pubsub.topics.topicMutex.RLock()
		topic, ok := pubsub.topics.topicsMap[topicID]
		pubsub.topics.topicMutex.RUnlock()
		if !ok {
			fmt.Println("Topic doesn't Exists")
			return false
		}

		pubsub.subscriptions.subscriptionMap[subName] = &Subscription{SubscriptionID: subName, pendingAck: make(map[int]*Message)}
		topic.Subscriptions = append(topic.Subscriptions, pubsub.subscriptions.subscriptionMap[subName])
		return true
	}
}

func (p *PubSub) DeleteSubscription(SubscriptionID string) {

	p.subscriptions.subMutex.Lock()
	defer p.subscriptions.subMutex.Unlock()

	_, ok := p.subscriptions.subscriptionMap[SubscriptionID]

	if ok {
		delete(p.topics.topicsMap, SubscriptionID)
		log.Printf("Subscription %q deleted \n", SubscriptionID)
	} else {
		log.Printf("SubscriptionID %q don't exist \n", SubscriptionID)
	}
}
