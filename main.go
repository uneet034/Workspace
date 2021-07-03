package main

import (
"fmt"

"time"
)

type Topic struct {
	message string
	id      int     //for identifying each topic
}

type Publisher struct {
	name string
}

type Broker struct {
	TopicBuffer chan Topic
	Subscribers map[int][]*Subscriber
}

type Subscriber struct {
	name   string
	topic  Topic
	buffer chan Topic
}

func (pub *Publisher) Publish(topic Topic, queue *Broker) bool {

	fmt.Println("Publisher is publishing the topic ", topic.message, "In process...")
	queue.TopicBuffer <- topic
	fmt.Println("\n Published the topic", topic.message, "To the message queue.")
	return true
}

func (pub *Publisher) SignalStop(queue *Broker) bool {

	return queue.SignalStop()
}

func (sub *Subscriber) Subscribe(queue *Broker) bool {

	fmt.Println("The Subscriber", sub.name, " is subscribing to the Topic ", sub.topic.message, "In process...")
	(*queue).Subscribers[sub.topic.id] = append((*queue).Subscribers[sub.topic.id], sub)
	fmt.Println("The Subscriber", sub.name, " is subscribed to the Topic ", sub.topic.message)
	return true
}

func (sub *Subscriber) ConsumeBuffer() bool {
	for topic := range (*sub).buffer {
		fmt.Println("Consumed ", topic.message, "from the Subscriber", sub.name)
	}
	fmt.Println("Subscriber ", sub.name, " Closed")
	return true
}

func (sub *Broker) NotifyConsumer() bool {
	for topic := range sub.TopicBuffer {
		subscribers := sub.Subscribers[topic.id]
		//fmt.Println(subscribers)
		for _, s := range subscribers {
			s.buffer <- topic
		}
	}
	return true
}

func (sub *Broker) SignalStop() bool {
	for _, v := range sub.Subscribers {
		for _, i := range v {
			close(i.buffer)
		}
	}
	return true
}

func main() {
	topics := []Topic{
		Topic{
			message: "first",
			id:      1,
		},
		Topic{
			message: "second",
			id:      2,
		},
		Topic{
			message: "third",
			id:      2,
		},
		Topic{
			message: "fourth",
			id:      2,
		},
		Topic{
			message: "fifth",
			id:      1,
		},
	}

	broker := Broker{
		TopicBuffer: make(chan Topic, 3),
		Subscribers: make(map[int][]*Subscriber),
	}

	publisher := Publisher{
		name: "first",
	}

	subscriber1 := Subscriber{
		name:   "S1",
		buffer: make(chan Topic),
		topic:  topics[0],
	}

	subscriber2 := Subscriber{
		name:   "S2",
		buffer: make(chan Topic),
		topic:  topics[1],
	}

	go subscriber1.ConsumeBuffer()
	go subscriber2.ConsumeBuffer()
	go broker.NotifyConsumer()

	subscriber1.Subscribe(&broker)
	subscriber2.Subscribe(&broker)

	for i := range topics {
		publisher.Publish(topics[i], &broker)
	}

	<-time.After(1 * time.Second)
	publisher.SignalStop(&broker)
	<-time.After(1 * time.Second)

}
