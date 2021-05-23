package pubsub

import (
	"log"
	"sync"
	"time"
)

type topic struct {
	topicId       string
	subscriptions []*subscription
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
	pendingAckMsg   map[int]*Message
	subscriber      *subscriber
	pendingMapMutex sync.RWMutex
}

// implicit implementation of iSubscription interface below

func (s *subscription) sendMessage(msg *Message) {

	s.pendingMapMutex.Lock()
	_, ok := s.pendingAckMsg[msg.messageId]
	if !ok {
		ok = true
		s.pendingAckMsg[msg.messageId] = msg
	}
	s.pendingMapMutex.Unlock()

	// Keep retry sending same message to subscriber each 5 second until we get the acknowledgement
	for ok {
		if s.subscriber != nil {
			go (*s.subscriber)(*msg)
		} else {
			log.Println("sendMessage()-> No subscriber function found, Maybe there is no one func has subscribed it : ", s.subscriptionID)
		}

		// 5 second time to check for Acknowledgement
		time.Sleep(retryTime)

		s.pendingMapMutex.RLock()
		_, ok = s.pendingAckMsg[msg.messageId]
		s.pendingMapMutex.RUnlock()

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
