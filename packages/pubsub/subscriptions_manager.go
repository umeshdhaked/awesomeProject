package pubsub

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type subscriptionTopics struct {
	subscriptionTopicMap      map[string]string
	subscriptionTopicMapMutex sync.RWMutex
}

//  Add subscription

func (p *PubSub) AddSubscription(topicID, subName string) (bool, error) {
	p.subscriptionTopics.subscriptionTopicMapMutex.Lock()
	defer p.subscriptionTopics.subscriptionTopicMapMutex.Unlock()

	_, ok := p.subscriptionTopics.subscriptionTopicMap[subName]
	if ok {
		fmt.Println("AddSubscription()-> subscription Already Exists",subName)
		return false, errors.New("subscription Already Exists")
	} else {

		p.topics.topicMutex.RLock()
		topicVar, ok := p.topics.topicsMap[topicID]
		p.topics.topicMutex.RUnlock()
		if !ok {
			fmt.Println("AddSubscription()-> topic doesn't Exists")
			return false, errors.New("topic doesn't Exists")
		}

		// adding subscription and topic entry in subscriptionTopicMap mutex
		p.subscriptionTopics.subscriptionTopicMap[subName] = topicID

		// adding subscription in topic's subscription Map
		topicVar.subscriptionsMutex.Lock()
		(*topicVar).subscriptions[subName] = &subscription{subscriptionID: subName, pendingAckMsg: make(map[int]*Message)}
		topicVar.subscriptionsMutex.Unlock()

		return true, nil
	}
}

func (p *PubSub) DeleteSubscription(SubscriptionID string)  (bool, error) {

	p.subscriptionTopics.subscriptionTopicMapMutex.Lock()
	defer p.subscriptionTopics.subscriptionTopicMapMutex.Unlock()

	topicId, ok := p.subscriptionTopics.subscriptionTopicMap[SubscriptionID]

	if ok {
		var (
			topicVar *topic
			ok1 bool
		)
		topicVar, ok1 = p.topics.topicsMap[topicId]

		if ok1 {
			// deleting subscription pointer from topic
			topicVar.subscriptionsMutex.Lock()
			delete(topicVar.subscriptions, SubscriptionID)
			topicVar.subscriptionsMutex.Unlock()
		}

		// deleting subscription entry from subscriptionTopicMap
		delete(p.subscriptionTopics.subscriptionTopicMap, SubscriptionID)

		log.Printf("DeleteSubscription()-> subscription %q deleted \n", SubscriptionID)

		return true,nil
	} else {
		p.subscriptionTopics.subscriptionTopicMapMutex.Unlock()
		log.Printf("DeleteSubscription()-> subscriptionID %q don't exist \n", SubscriptionID)
		return false, errors.New("subscriptionID don't exist")
	}
}
