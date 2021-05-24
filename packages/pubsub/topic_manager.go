package pubsub

import (
	"errors"
	"fmt"
	"log"
)
import "sync"

type topics struct {
	topicsMap  map[string]*topic
	topicMutex sync.RWMutex
}

func (p *PubSub) CreateTopic(topicName string) (bool, error) {
	p.topics.topicMutex.Lock()
	defer p.topics.topicMutex.Unlock()
	_, ok := p.topics.topicsMap[topicName]

	if ok {
		fmt.Println("topic Already Exists", topicName)
		return false, errors.New("topic already exists")
	} else {
		p.topics.topicsMap[topicName] = &topic{topicId: topicName, subscriptions: make(map[string]*subscription)}
		return true, nil
	}
}

func (p *PubSub) DeleteTopic(TopicID string) (bool, error) {

	p.topics.topicMutex.Lock()
	defer p.topics.topicMutex.Unlock()

	topicVar, ok := p.topics.topicsMap[TopicID]
	if ok {

		// first removing all subscriptions of that Topic from subscriptionTopicsMap
		topicVar.subscriptionsMutex.RLock()
		p.subscriptionTopics.subscriptionTopicMapMutex.Lock()
		for key,_ := range topicVar.subscriptions {
			delete(p.subscriptionTopics.subscriptionTopicMap,key)
		}
		p.subscriptionTopics.subscriptionTopicMapMutex.Lock()
		topicVar.subscriptionsMutex.RUnlock()

		//then deleting topic from topicsMap
		delete(p.topics.topicsMap, TopicID)

		log.Printf("topic %q deleted \n", TopicID)
		return true, nil
	} else {
		log.Printf("TopicID %q don't exist \n", TopicID)
		return false,errors.New("topicID don't exist")
	}

}
