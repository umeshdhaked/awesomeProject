package pubsub

import "fmt"
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
	fmt.Println("DeleteTopic yet to be implemented")
}
