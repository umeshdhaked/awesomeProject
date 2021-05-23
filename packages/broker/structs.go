package broker


type Topic struct {
	topicId string
	subscriptions []*Subscription
}

type Subscription struct {
	subscriptionId string
}
