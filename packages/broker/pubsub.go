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

	var subscription = subscriptions.subscriptionMap[subscriptionID]

	subscriber := Subscriber(subscriberFunc)

	subscription.addSubscriber(&subscriber)

}

func Unsubscribe(subId string) {
	var subscription = subscriptions.subscriptionMap[subId]
	subscription.removeSubscriber()
}

var flag = true
var ch = make(chan Message, 10)

func Publish(topicId, message string) {

	messageIdTracker++
	messageObj := Message{messageIdTracker, topicId, message}

	if flag {
		go pushMessage(ch)
		flag = false
	}
	ch <- messageObj
}

func pushMessage(ch chan Message) {

	for {
		msg := <-ch
		var topic = topics.topicsMap[msg.TopicId]

		var subsObjs = topic.Subscriptions

		for _, val := range subsObjs {
			go val.sendMessage(&msg)
		}
	}

}
