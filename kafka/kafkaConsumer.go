package kafka

import (
	"encoding/json"
	"fmt"
	"mailSender/service"

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
		"auto.offset.reset": "earliest",         // Adjust as needed based on your requirements
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

type MessageContent struct {
	CarPlate         string  `json:"carPlate"`
	Address          string  `json:"address"`
	Date             string  `json:"date"`
	Violation        string  `json:"violation"`
	CarType          string  `json:"carType"`
	CarColor         string  `json:"carColor"`
	CarBrand         string  `json:"carBrand"`
	VehicleOwnerName string  `json:"vehicleOwnerName"`
	VehicleOwnerCPF  string  `json:"veiculeOwneCPF"`
	Speed            float64 `json:"speed"`
	MaxSpeed         int     `json:"maxSpeed"`
	FinePrice        float64 `json:"finePrice"`
	Sex              string  `json:"sex"`
	Age              int     `json:"age"`
}

func (kc *KafkaClient) ConsumeMessages(topic string) {
	kc.Consumer.SubscribeTopics([]string{topic}, nil)

	fmt.Printf("Consumer subscribed to topic: %s\n", topic)

	for {
		msg, err := kc.Consumer.ReadMessage(-1)
		if err == nil {
			// Decodifica o valor da mensagem
			var messageContent MessageContent
			err := json.Unmarshal(msg.Value, &messageContent)
			if err != nil {
				fmt.Printf("Error decoding message value: %v\n", err)
				continue
			}

			// Convert messageContent to service.MessageContent
			serviceMessageContent := service.MessageContent{
				CarPlate:         messageContent.CarPlate,
				Address:          messageContent.Address,
				Date:             messageContent.Date,
				Violation:        messageContent.Violation,
				CarType:          messageContent.CarType,
				CarColor:         messageContent.CarColor,
				CarBrand:         messageContent.CarBrand,
				VehicleOwnerName: messageContent.VehicleOwnerName,
				VehicleOwnerCPF:  messageContent.VehicleOwnerCPF,
				Speed:            messageContent.Speed,
				MaxSpeed:         messageContent.MaxSpeed,
				FinePrice:        messageContent.FinePrice,
				Sex:              messageContent.Sex,
				Age:              messageContent.Age,
			}

			// Show infraction details
			service.ShowInfraction(serviceMessageContent)

			// Get email
			email := service.GetEmail(serviceMessageContent)
			fmt.Println("Email: ", email)

			// Validate CPF and send email
			if service.ValidateCPF(serviceMessageContent) {
				// Envia o email
				err := service.SetupEmail(serviceMessageContent, email)
				if err != nil {
					fmt.Println(err)
				}
			}

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
