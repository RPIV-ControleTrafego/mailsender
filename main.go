//go:build ignore
// +build ignore

package main

import (

	"log"
	
	"mailSender/kafka"    

	
)



func main() {
	// Crie uma configuração Kafka
	kafkaConfig := kafka.NewKafkaConfiguration()

	// Crie um cliente Kafka
	kafkaClient, err := kafka.NewKafkaClient(kafkaConfig)
	if err != nil {
		log.Fatalf("Error initializing Kafka client: %v", err)
		return
	}
	defer kafkaClient.Close()

	// Chame a função ConsumeMessages para consumir mensagens do Kafka
	go kafkaClient.ConsumeMessages("infraction-topic")

	
	// Faça o que mais for necessário na sua aplicação...

	// Agora, o programa permanecerá em execução para consumir mensagens Kafka
	select {}
}



