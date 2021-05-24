package pubsub

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type subscriptionTopics struct {
	subscriptionTopicMap      sync.Map
	subscriptionTopicMapMutex sync.RWMutex
}

//  Add subscription

func (p *PubSub) AddSubscription(topicID, subName string) (bool, error) {

	_, ok := p.subscriptionTopics.subscriptionTopicMap.Load(subName)
	if ok {
		fmt.Println("AddSubscription()-> subscription Already Exists", subName)
		return false, errors.New("subscription Already Exists")
	} else {

		topicVar, ok := p.topics.topicsMap.Load(topicID)
		if !ok {
			fmt.Println("AddSubscription()-> topic doesn't Exists")
			return false, errors.New("topic doesn't Exists")
		}

		// adding subscription and topic entry in subscriptionTopicMap mutex
		p.subscriptionTopics.subscriptionTopicMap.Store(subName, topicID)

		// adding subscription in topic's subscription Map
		topicVar.(*topic).subscriptions.Store(subName, &subscription{subscriptionID: subName, pendingAckMsg: sync.Map{}})

		return true, nil
	}
}

func (p *PubSub) DeleteSubscription(SubscriptionID string) (bool, error) {

	topicId, ok := p.subscriptionTopics.subscriptionTopicMap.Load(SubscriptionID)

	if ok {
		topicVar, ok1 := p.topics.topicsMap.Load(topicId.(string))
		if ok1 {
			// removing subscriber before deleting subscription
			p.UnSubscribe(SubscriptionID)
			// deleting subscription pointer from topic
			topicVar.(*topic).subscriptions.Delete(SubscriptionID)
		}

		// deleting subscription entry from subscriptionTopicMap
		p.subscriptionTopics.subscriptionTopicMap.Delete(SubscriptionID)

		log.Printf("DeleteSubscription()-> subscription %q deleted \n", SubscriptionID)

		return true, nil
	} else {
		log.Printf("DeleteSubscription()-> subscriptionID %q don't exist \n", SubscriptionID)
		return false, errors.New("subscriptionID don't exist")
	}
}
