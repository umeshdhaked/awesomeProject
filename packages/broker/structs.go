package broker

import "log"

type Topic struct {
	TopicId       string
	Subscriptions []*Subscription
}

type Message struct {
	MessageId int
	TopicId   string
	Data      string
}

type Subscriber func(message Message)

type ISubscription interface {
	sendMessage(msg *Message)
	addSubscriber(subscriber *Subscriber)
	removeSubscriber()
}

type Subscription struct {
	SubscriptionID string
	Subscriber     *Subscriber
}

func (s *Subscription) sendMessage(msg *Message) {

	if s.Subscriber != nil {
		(*s.Subscriber)(*msg)
	} else {
		log.Println("No subscription function found, Maybe there is no subscriber for subscription : ", s.SubscriptionID)
	}

}

func (s *Subscription) addSubscriber(subscriber *Subscriber) {
	if s.Subscriber == nil {
		s.Subscriber = subscriber
	} else {
		log.Println("Subscriber already added, request rejected!")
		//s.Subscriber = subscriber
	}
}

func (s *Subscription) removeSubscriber() {
	if s.Subscriber == nil {
		log.Println("No subscriber found")
	} else {
		s.Subscriber = nil
		log.Println("Subscriber removed")
	}
}
