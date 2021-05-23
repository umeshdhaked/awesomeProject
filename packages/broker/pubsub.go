package broker

type PubSub interface {
	CreateTopic()

	AddSubscription()

	Publish()

	Subscribe()

	Ack()
}

var messageIdTracker = 0

func GetAllTopicsAndSubscriptions() (map[string]*Topic, map[string]*Subscription) {
	return topics.topicsMap, subscriptions.subscriptionMap
}

func Subscribe(subscriptionID string, subscriberFunc func(msg Message)) {

	var subscription *Subscription = subscriptions.subscriptionMap[subscriptionID]

	subscriber := Subscriber(subscriberFunc)

	subscription.addSubscriber(&subscriber)

}

func Unsubscribe(subId string) {
	var subscription *Subscription = subscriptions.subscriptionMap[subId]
	subscription.removeSubscriber()
}

func Publish(topicId, message string) {
	messageIdTracker++
	messageObj := Message{messageIdTracker, message}

	pushMessage(topicId, &messageObj)

}

func pushMessage(topicId string, msg *Message) {

	var topic = topics.topicsMap[topicId]

	var subsObjs = topic.Subscriptions

	for _, val := range subsObjs {
		val.sendMessage(msg)
	}

}