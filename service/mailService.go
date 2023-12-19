package service

import (
	"fmt"
	"log"
	"mailSender/database"
	"net/smtp"
)

// Observer interface
type Observer interface {
	Update(message MessageContent)
}

// EmailService struct
type EmailService struct {
	Observers []Observer
}

// AddObserver method to add an observer
func (es *EmailService) AddObserver(observer Observer) {
	es.Observers = append(es.Observers, observer)
}

// NotifyObservers method to notify all observers
func (es *EmailService) NotifyObservers(message MessageContent) {
	for _, observer := range es.Observers {
		observer.Update(message)
	}
}

// MessageContent struct
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

// ShowInfraction displays the details of an infraction and notifies a list of observers.
func ShowInfraction(message MessageContent, observers ...Observer) {
	log.Println("Infraction details:")
	log.Printf("   Car Plate: %s\n", message.CarPlate)
	log.Printf("   Violation: %s\n", message.Violation)
	log.Printf("   Owner Name: %s\n", message.VehicleOwnerName)
	log.Printf("   Owner CPF: %s\n", message.VehicleOwnerCPF)
// //fine price
// 	log.Printf("   Fine Price: %.2f\n", message.FinePrice)

	// Notify observers
	for _, observer := range observers {
		observer.Update(message)
	}
}

var proxy = database.NewMongoDBProxy()

func ValidateCPF(message MessageContent) bool {
	email, err := proxy.GetEmailByCPF(message.VehicleOwnerCPF)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if email == "" {
		fmt.Println("CPF not found")
		return false
	}

	fmt.Println("Email: ", email)

	return true
}

func GetEmail(message MessageContent) string {
	email, err := proxy.GetEmailByCPF(message.VehicleOwnerCPF)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if email == "andremendes0113@gmail.com" {
		return "andremiranda.aluno@unipampa.edu.br"
	}

	return email
}

func SendMail(from, password, host string, port int, to []string, message []byte) error {
	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func SetupEmail(message MessageContent, to string, observers ...Observer) error {
	from := "andremendes0113@gmail.com"
	password := "0294BE2A8D274556D9584EE59D90DEFC0AB6" // Sua senha do Elastic Email

	subject := "Infraction"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		"Placa: " + message.CarPlate + "\n" +
		"Prezado, informamos que uma multa em seu CPF foi emetida em nossos sistemas, para poder pagá-la acesse o site TodayTrafficSystem" + "\n" +
		"Violação: " + message.Violation + "\n" +
		"Nome: " + message.VehicleOwnerName + "\n" +
		"CPF: " + message.VehicleOwnerCPF + "\n" +
		"**Valor da Multa:** " + fmt.Sprintf("%.2f", message.FinePrice) + "\n" + // Include only once
		"Endereço: " + message.Address + "\n" +
		"Data: " + message.Date + "\n" +
		"Tipo de Carro: " + message.CarType + "\n" +
		"Cor do Carro: " + message.CarColor + "\n" +
		"Marca do Carro: " + message.CarBrand + "\n" +
		"Velocidade: " + fmt.Sprintf("%.2f", message.Speed) + "\n" +
		"Velocidade Máxima Permitida: " + fmt.Sprintf("%d", message.MaxSpeed) + "\n" +
		"\n"
	err := SendMail(from, password, "smtp.elasticemail.com", 2525, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("Erro ao enviar o email:", err)
	} else {
		fmt.Println("Email enviado com sucesso para:", to)
	}

	// Notify observers
	for _, observer := range observers {
		observer.Update(message)
	}

	return nil
}

// EmailObserver struct
type EmailObserver struct{}

// Update method for EmailObserver
func (eo *EmailObserver) Update(message MessageContent) {
	fmt.Println("Observer received update:")
	fmt.Printf("   Car Plate: %s\n", message.CarPlate)
	fmt.Printf("   Violation: %s\n", message.Violation)
	fmt.Printf("   Owner Name: %s\n", message.VehicleOwnerName)
	fmt.Printf("   Owner CPF: %s\n", message.VehicleOwnerCPF)
}
