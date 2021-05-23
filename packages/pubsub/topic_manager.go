package pubsub

import "fmt"
import "sync"

type Topics struct {
	topicsMap map[string]*Topic
	topicMutex sync.RWMutex
}


//  Create topic
var topics Topics = Topics{topicsMap: make(map[string]*Topic)}

func CreateTopic(topicName string) bool {
	topics.topicMutex.Lock()
	defer topics.topicMutex.Unlock()
	val, ok := topics.topicsMap[topicName]

	if ok {
		fmt.Println("Topic Already Exists", val)
		return false
	} else {
		topics.topicsMap[topicName] = &Topic{TopicId: topicName}
		return true
	}
}