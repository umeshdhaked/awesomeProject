package pubsub

import (
	"log"
	"sync"
	"time"
)

type topic struct {
	topicId            string
	subscriptions      sync.Map
}

type Message struct {
	messageId int
	topicId   string
	data      string
}

func (m Message) MessageId() int {
	return m.messageId
}

func (m Message) TopicId() string {
	return m.topicId
}

func (m Message) Data() string {
	return m.data
}

type subscriber func(message Message)

type iSubscription interface {
	sendMessage(msg *Message)
	addSubscriber(subscriber *subscriber)
	removeSubscriber()
}

type subscription struct {
	subscriptionID  string
	pendingAckMsg   sync.Map
	subscriber      *subscriber
}

// implicit implementation of iSubscription interface below

func (s *subscription) sendMessage(msg *Message) {

	_, ok := s.pendingAckMsg.Load(msg.messageId)
	if !ok {
		ok = true
		s.pendingAckMsg.Store(msg.messageId, msg)
	}

	// Keep retry sending same message to subscriber each 5 second until we get the acknowledgement
	for i := 0; i < retryMaxCount && ok; i++ {
		if s.subscriber != nil {
			go (*s.subscriber)(*msg)
		} else {
			log.Println("sendMessage()-> No subscriber found (Stopping retrySend also, message dropped for this subscriber) for subscription : ", s.subscriptionID)
			s.pendingAckMsg.Delete(msg.messageId)
			ok = false
			continue
		}

		// 15 second default time to check for Acknowledgement
		time.Sleep(retryTime)

		_, ok = s.pendingAckMsg.Load(msg.messageId)

		if ok {
			log.Printf("sendMessage()-> Resending the message with id %v \n", msg.messageId)
		}
	}

}

func (s *subscription) addSubscriber(subscriber *subscriber) {
	if s.subscriber == nil {
		s.subscriber = subscriber
	} else {
		log.Println("addSubscriber()-> subscriber already added, request rejected!")
		//s.subscriber = subscriber
	}
}

func (s *subscription) removeSubscriber() {
	if s.subscriber == nil {
		log.Println("removeSubscriber()-> No subscriber found")
	} else {
		s.subscriber = nil
		log.Printf("removeSubscriber()-> subscriber removed from subscription %q ", s.subscriptionID)
	}
}
