package pubsub

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type subscriptions struct {
	subscriptionMap map[string]*subscription
	subMutex        sync.RWMutex
}

//  Add subscription

func (p *PubSub) AddSubscription(topicID, subName string) (bool, error) {
	p.subscriptions.subMutex.Lock()
	defer p.subscriptions.subMutex.Unlock()

	val, ok := p.subscriptions.subscriptionMap[subName]
	if ok {
		fmt.Println("AddSubscription()-> subscription Already Exists", val)
		return false, errors.New("subscription Already Exists")
	} else {
		p.topics.topicMutex.RLock()
		topic, ok := p.topics.topicsMap[topicID]
		p.topics.topicMutex.RUnlock()
		if !ok {
			fmt.Println("AddSubscription()-> topic doesn't Exists")
			return false, errors.New("topic doesn't Exists")
		}

		p.subscriptions.subscriptionMap[subName] = &subscription{subscriptionID: subName, pendingAckMsg: make(map[int]*Message)}
		topic.subscriptions = append(topic.subscriptions, p.subscriptions.subscriptionMap[subName])
		return true, nil
	}
}

func (p *PubSub) DeleteSubscription(SubscriptionID string)  (bool, error) {

	p.subscriptions.subMutex.Lock()
	defer p.subscriptions.subMutex.Unlock()

	_, ok := p.subscriptions.subscriptionMap[SubscriptionID]

	if ok {
		delete(p.topics.topicsMap, SubscriptionID)
		log.Printf("subscription %q deleted \n", SubscriptionID)

		return true,nil
	} else {
		log.Printf("DeleteSubscription()-> subscriptionID %q don't exist \n", SubscriptionID)
		return false, errors.New("subscriptionID don't exist")
	}
}
