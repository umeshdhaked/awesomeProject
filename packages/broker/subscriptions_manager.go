package broker

import (
	"fmt"
	"sync"
)

type Subscriptions struct {
	subscriptionMap map[string]*Subscription
	subMutex sync.RWMutex
}

//  Add subscription
var subscriptions Subscriptions = Subscriptions{subscriptionMap: make(map[string]*Subscription)}

func AddSubscription(topicID string, subName string) bool {
	subscriptions.subMutex.Lock()
	defer subscriptions.subMutex.Unlock()

	val, ok := subscriptions.subscriptionMap[subName]
	if ok {
		fmt.Println("Subscription Already Exists", val)
		return false
	} else {
		topics.topicMutex.RLock()
		topic, ok := topics.topicsMap[topicID]
		topics.topicMutex.RUnlock()
		if !ok {
			fmt.Println("Topic doesn't Exists")
			return false
		}

		subscriptions.subscriptionMap[subName] = &Subscription{SubscriptionID: subName}
		topic.Subscriptions = append(topic.Subscriptions, subscriptions.subscriptionMap[subName] )
		return true
	}
}