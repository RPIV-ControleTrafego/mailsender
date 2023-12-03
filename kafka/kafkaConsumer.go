package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const PartitionAny = -1

type KafkaClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

type KafkaConfiguration struct {
	BootstrapServers string
}

func NewKafkaClient(config KafkaConfiguration) (*KafkaClient, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	})

	if err != nil {
		return nil, err
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          "infraction-topic", // Specify your consumer group ID
		"auto.offset.reset": "earliest",      // Adjust as needed based on your requirements
	})

	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		Producer: p,
		Consumer: c,
	}, nil
}

func (kc *KafkaClient) Close() {
	kc.Producer.Close()
	kc.Consumer.Close()
}

func (kc *KafkaClient) SendMessage(topic string, message string) error {
	deliveryChan := make(chan kafka.Event)

	err := kc.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	fmt.Printf("Message delivered to topic %s (partition %d) at offset %d\n", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	return nil
}

func (kc *KafkaClient) ConsumeMessages(topic string) {
    kc.Consumer.SubscribeTopics([]string{topic}, nil)

    fmt.Printf("Consumer subscribed to topic: %s\n", topic)

    for {
        msg, err := kc.Consumer.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Received message: %s\n", string(msg.Value))
        } else {
            fmt.Printf("Error receiving message: %v\n", err)
        }
    }
}

func NewKafkaConfiguration() KafkaConfiguration {
	return KafkaConfiguration{
		BootstrapServers: "127.0.0.1:9092",
	}
}
