package broker

import "log"

type Topic struct {
	TopicId       string
	Subscriptions []*Subscription
}

type Subscription struct {
	SubscriptionID string
	Subscriber     *Subscriber
}

type Message struct {
	MessageId int
	Data      string
}

type Subscriber func(message Message)


func (s Subscription) sendMessage(msg Message) {

	if s.Subscriber != nil {
		(*s.Subscriber)(msg)
	} else {
		log.Println("No subscription function found, Maybe there is no subscriber for subscription : ", s.SubscriptionID)
	}

}
