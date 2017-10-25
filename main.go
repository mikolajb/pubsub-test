package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

const nbOfMessages = 1000000

func main() {
	fmt.Println("START")
	ctx, cancel := context.WithCancel(context.Background())

	// Sets your Google Cloud Platform project ID.
	projectID := "test-project"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new topic.
	topicName := "my-new-topic"

	// Creates the new topic.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	fmt.Printf("Topic %v created.\n", topic)

	done := make(chan struct{})

	go func() {
		for i := uint64(0); i < nbOfMessages; i++ {
			b := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(b, i)
			topic.Publish(ctx, &pubsub.Message{Data: b})
		}

		fmt.Println("Producer is done")
		<-done
		cancel()
	}()

	sub, err := client.CreateSubscription(ctx, "test-sub", pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		log.Fatalf("Cannot create a sub: %s", err.Error())
	}

	m := &sync.Mutex{}
	set := make(map[uint64]bool)
	c := 0

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		m.Lock()
		i, read := binary.Uvarint(msg.Data)
		if set[i] {
			fmt.Printf("duplicate: %8d\n", i)
			return
		}
		set[i] = true
		c++
		if c == nbOfMessages-1 {
			close(done)
		}
		m.Unlock()
		if c%10000 == 0 {
			fmt.Printf("%8d %8d %d\n", c, i, read)
		}
		msg.Ack()
	})
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	fmt.Println("It read: ", c, "messages")
}
