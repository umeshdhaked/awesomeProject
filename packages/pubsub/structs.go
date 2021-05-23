package pubsub

import (
	"log"
	"sync"
	"time"
)

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
	SubscriptionID  string
	pendingAck      map[int]*Message
	Subscriber      *Subscriber
	pendingMapMutex sync.RWMutex
}

// implicit implementation of ISubscription interface below

func (s *Subscription) sendMessage(msg *Message) {

	s.pendingMapMutex.Lock()
	_, ok := s.pendingAck[msg.MessageId]
	if !ok {
		ok = true
		s.pendingAck[msg.MessageId] = msg
	}
	s.pendingMapMutex.Unlock()

	// Keep retry sending same message to subscriber each 5 second until we get the acknowledgement
	for ok {
		if s.Subscriber != nil {
			go (*s.Subscriber)(*msg)
		} else {
			log.Println("No subscription function found, Maybe there is no subscriber for subscription : ", s.SubscriptionID)
		}

		// 5 second time to check for Acknowledgement

		time.Sleep(time.Second * 5)

		s.pendingMapMutex.RLock()
		_, ok = s.pendingAck[msg.MessageId]
		s.pendingMapMutex.RUnlock()

		if ok {
			log.Printf("MessageID %q not acknowledgd by subscriber. resending the message \n", string(msg.MessageId))
		}
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
		log.Printf("Subscriber removed from Subscription %q ", s.SubscriptionID)
	}
}
