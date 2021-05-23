package pubsub

import (
	"fmt"
	"log"
)
import "sync"

type Topics struct {
	topicsMap  map[string]*Topic
	topicMutex sync.RWMutex
}

func (p *PubSub) CreateTopic(topicName string) bool {
	pubsub.topics.topicMutex.Lock()
	defer pubsub.topics.topicMutex.Unlock()
	val, ok := pubsub.topics.topicsMap[topicName]

	if ok {
		fmt.Println("Topic Already Exists", val)
		return false
	} else {
		pubsub.topics.topicsMap[topicName] = &Topic{TopicId: topicName}
		return true
	}
}

func (p *PubSub) DeleteTopic(TopicID string) {

	p.topics.topicMutex.Lock()
	defer p.topics.topicMutex.Unlock()

	_, ok := p.topics.topicsMap[TopicID]
	if ok {
		delete(p.topics.topicsMap, TopicID)
		log.Printf("Topic %q deleted \n", TopicID)
	} else {
		log.Printf("TopicID %q don't exist \n", TopicID)
	}

}
