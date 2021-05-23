package broker

type PubSub interface {
	CreateTopic()

	AddSubscription()

	Publish()

	Subscribe()

	Ack()
}

func GetAll() (map[string]*Topic, map[string]*Subscription){
	return topics,subscriptions
}

