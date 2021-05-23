package broker


type Topic struct {
	TopicId       string
	Subscriptions []*Subscription
}

type Subscription struct {
	SubscriptionId string
}
