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
   **e**. Library will keep resending a message every 15 second to a subscriber until it gets the acknowledgement from that subscriber for message
   or max retryCount(which is 5 by default) exceeds. <br />
   

4. Running embedded test runner <br/>
   To run internal test runner main package, Download lib and unzip it, then navigate to runner directory `cd runner` <br/>
   then run following command _`go run main.go`_ <br/>

   It will create topics, {topicID : topic1} and {topicID : topic2} <br/>
   It will create three subscriptions, {topicID: topic1, subscriptionID: topic1_sub1}, {topicID: topic1, subscriptionID: topic1_sub2} and {topicID: topic2, subscriptionID: topic2_sub1} <br/>

   It will create three subscribers, {SubscriptionID: topic1_sub1, subsFunc: SubscriberTypeA}, {SubscriptionID: topic1_sub2, subsFunc: SubscriberTypeB} and {SubscriptionID: topic2_sub1, subsFunc: SubscriberTypeC} <br/>\

   then will publish random messages to "topic1" and "topic2" which will propagate each message to all the subscriptions of topics, and you will be able to see output on your screen. <br/>
   Then, I am deleting the "topic1" and publish a message after deleting topic which returns an error.
   
5. Assumptions made <br/>
   **a.** Deleting a topic will also delete its subscriptions and subscribers of those subscriptions as well.<br/>
   **b.** If message is pushed and topic deleted/subscription deleted/unsubscribed before ACK then message will not be sent again. It will show an error in logs.<br/>
   **c.** If subscription or subscriber is deleted for sending ACK then error message will be shown.<br/>