package subscribers

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"log"
)

func SubscriberTypeA(msg pubsub.Message) {
	fmt.Println("SubscriberTypeA:", msg.Data)
}

func SubscriberTypeB(msg pubsub.Message) {
	log.Println("SubscriberTypeB:", msg.Data)
}

func SubscriberTypeC(msg pubsub.Message) {
	fmt.Println("SubscriberTypeC:", msg.Data)
}

func SubscriberTypeD(msg pubsub.Message) {
	fmt.Println("SubscriberTypeD:", msg.Data)
}

func SubscriberTypeE(msg pubsub.Message) {
	fmt.Println("SubscriberTypeE:", msg.Data)
}

func SubscriberTypeF(msg pubsub.Message) {
	fmt.Println("SubscriberTypeF:", msg.Data)
}

func SubscriberTypeG(msg pubsub.Message) {
	fmt.Println("SubscriberTypeG:", msg.Data)
}
