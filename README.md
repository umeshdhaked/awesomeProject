# Golang PubSub system
1. This is a PubSub library which can be embedded in any system where you want to publish data to a **_topic_** 
and consume at many **_subscriptions_**

2. Library exposes following apis to use it. <br />
      **CreateTopic(topicName string)** `create a topic if topic not exist`
   
      **DeleteTopic(TopicID string)** `delete an existing topic`
   
      **AddSubscription(topicID, subName string)** `create new subscription for an existing topic`
   
      **DeleteSubscription(SubscriptionID string)** `delete an existing subscription`
   
      **Subscribe(subscriptionID string, subscriberFunc func(msg Message))** `when message is published on a topic of this subscription then provided SubscriberFunc will be executed`
   
      **UnSubscribe(subId string)** `after unsubscribe, SubscriptionFunc of this subscription will not be executed`

      **Publish(topicId, message string)** `this will add the message to topic and push it to all the subscriptions of topic then executes the subscriberFunc`

      **Ack(msgId int, subId string)** `subscriber will have to acknowledge the received message by calling this method`


3. Some important points to keep in mind. <br />
   **a**. One-to-many relation of topic with subscriptions <br />
   **b**. One-to-one relation of subscription with subscriber. <br />
   **c**. Publish-Subscribe system puts the messages from Publish api to a buffered channel of default bufferSize 50. <br />
   **d**. Internally library waits 15 second for the pushed message to be acknowledged by subscriber before resending it again. <br />
   **e**. Library will keep resending a message every 15 second to a subscriber until it gets the acknowledgement from that subscriber for message. <br />
   

4. Running embedded test runner <br/>
   To run internal test runner main package, Download the zip.
   then run from base directory run _`go run runner/main.go`_ <br/>
   
   it will create one topic "topic1" and two subscription "sub1" and "sub2"
   then will do a 100 of publish to "topic1" which will propagate each message to both subscriptions, and you will be able to see output on your screen