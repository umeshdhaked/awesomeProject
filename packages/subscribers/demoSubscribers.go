package subscribers

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"time"
)

func SubscriberTypeA(msg pubsub.Message) {
	ps := pubsub.GetPubSub()

	fmt.Println("SubscriberTypeA:", msg.MessageId, msg.Data, msg.TopicId)

	ps.Ack(msg.MessageId, "sub1")
}

func SubscriberTypeB(msg pubsub.Message) {
	ps := pubsub.GetPubSub()

	fmt.Println("SubscriberTypeB:", msg.MessageId, msg.Data, msg.TopicId)
	time.Sleep(time.Second * 8)

	ps.Ack(msg.MessageId, "sub2")

}

func SubscriberTypeC(msg pubsub.Message) {
	fmt.Println("SubscriberTypeC:", msg.MessageId, msg.Data)
}

func SubscriberTypeD(msg pubsub.Message) {
	fmt.Println("SubscriberTypeD:", msg.MessageId, msg.Data)
}

func SubscriberTypeE(msg pubsub.Message) {
	fmt.Println("SubscriberTypeE:", msg.MessageId, msg.Data)
}

func SubscriberTypeF(msg pubsub.Message) {
	fmt.Println("SubscriberTypeF:", msg.MessageId, msg.Data)
}

func SubscriberTypeG(msg pubsub.Message) {
	fmt.Println("SubscriberTypeG:", msg.MessageId, msg.Data)
}
