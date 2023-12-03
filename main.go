//go:build ignore
// +build ignore

package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/gofiber/fiber/v2"

	// "mailSender/database" // update import path
	"mailSender/kafka"    // update import path

	"github.com/arsmn/fiber-swagger/v2"
)

// @title Email Service API
// @version 1.0
// @description API for sending emails using Fiber
// @host localhost:3000

// @BasePath /

type EmailAndCPF struct {
	TrafficViolation TrafficViolation `json:"trafficViolation"`
	CPF              string           `json:"cpf"`
}

// Método para converter TrafficViolation em Email
func (tv TrafficViolation) ToEmail() Email {
	return Email{
		To:      "to@example.com",  // Substitua pelo destino desejado
		Subject: "Traffic Violation",
		Body: fmt.Sprintf(
			"Traffic Violation Details:\n\n"+
				"Car Plate: %s\n"+
				"Address: %s\n"+
				// Adicione outros campos conforme necessário
				"Speed: %.2f\n"+
				"Max Speed: %d\n"+
				"Fine Price: %.2f\n"+
				// Adicione outros campos conforme necessário
				"Vehicle Owner Name: %s\n"+
				"Vehicle Owner CPF: %s\n",
			tv.CarPlate, tv.Address, tv.Speed, tv.MaxSpeed, tv.FinePrice, tv.VehicleOwnerName, tv.VehicleOwnerCPF),
	}
}

// func main() {
// 	// Criação e inicialização do Service Locator
// 	serviceLocator := NewEmailServiceLocator()
// 	serviceLocator.Initialize(
// 		"andremendes0113@gmail.com",
// 		"0294BE2A8D274556D9584EE59D90DEFC0AB6",
// 		"smtp.elasticemail.com",
// 		2525,
// 	)

// 	// Inicialização do consumidor Kafka em segundo plano
// 	kafkaConfig := kafka.NewKafkaConfiguration()
// 	kafkaClient, err := kafka.NewKafkaClient(kafkaConfig)
// 	if err != nil {
// 		fmt.Printf("Error initializing Kafka client: %v\n", err)
// 		return
// 	}
// 	defer kafkaClient.Close()
// 	go kafkaClient.ConsumeMessages("infraction-topic")
// // Loop de leitura de mensagens do Kafka
// for {
// 	msg, err := kafkaClient.Consumer.ReadMessage(-1)
// 	if err == nil {
// 		var emailAndCPF EmailAndCPF
// 		if err := json.Unmarshal(msg.Value, &emailAndCPF); err != nil {
// 			log.Printf("Error decoding Kafka message: %v", err)
// 			continue
// 		}

// 		// Verifique se o CPF existe no banco de dados antes de enviar o e-mail
// 		exists, err := database.CheckCPFsExistInDB(emailAndCPF.CPF)
// 		if err != nil {
// 			log.Printf("Error checking CPF existence: %v", err)
// 			continue
// 		}

// 		if !exists {
// 			log.Printf("CPF %s does not exist in the database.", emailAndCPF.CPF)
// 			continue
// 		}

// 		// Converta TrafficViolation em Email e envie o e-mail
// 		email := emailAndCPF.TrafficViolation.ToEmail()
// 		err = serviceLocator.GetService().SendEmail(email)
// 		if err != nil {
// 			log.Printf("Error sending email to %s: %v", email.To, err)
// 			continue
// 		}

// 		log.Printf("Email sent successfully to %s", email.To)
// 	} else {
// 		log.Printf("Error reading Kafka message: %v", err)
// 	}
// }
// }


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




func processMessage(msg string) {
    // Adicione aqui o código para processar a mensagem recebida
    log.Printf("Processing message: %s", msg)
    // ... outras operações ...
}


// @Summary Send an email
// @Description Send an email using the provided details.
// @Tags email
// @Accept json
// @Produce json
// @Param email body Email true "Email details"
// @Success 200 {string} string "Email sent successfully."
// @Router /email [post]
func setupApp(sl *EmailServiceLocator) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	
	
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/email", func(c *fiber.Ctx) error {
		return c.SendString(sl.GetView().SuccessMessage())
	})

	app.Post("/email", func(c *fiber.Ctx) error {
		var email Email
		if err := c.BodyParser(&email); err != nil {
			return err
		}

		err := sl.GetService().SendEmail(email)
		if err != nil {
			return c.Status(500).SendString(sl.GetView().ErrorMessage())
		}

		return c.SendString(sl.GetView().SuccessMessage())
	})

	return app
}

// Model
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EmailServiceLocator é o Service Locator responsável por fornecer instâncias dos serviços necessários.
type EmailServiceLocator struct {
	Controller *EmailController
	View       *EmailView
	Service    *EmailService
}

func NewEmailServiceLocator() *EmailServiceLocator {
	return &EmailServiceLocator{}
}

func (sl *EmailServiceLocator) Initialize(from, password, host string, port int) {
	sl.Controller = NewEmailController(from, password, host, port)
	sl.View = NewEmailView()
	sl.Service = NewEmailService(sl.Controller, sl.View)
}

func (sl *EmailServiceLocator) GetController() *EmailController {
	return sl.Controller
}

func (sl *EmailServiceLocator) GetView() *EmailView {
	return sl.View
}

func (sl *EmailServiceLocator) GetService() *EmailService {
	return sl.Service
}

type EmailController struct {
	From     string
	Password string
	Host     string
	Port     int
}

func NewEmailController(from, password, host string, port int) *EmailController {
	return &EmailController{
		From:     from,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

type EmailView struct{}

func NewEmailView() *EmailView {
	return &EmailView{}
}

func (v *EmailView) SuccessMessage() string {
	return "Email sent successfully."
}

func (v *EmailView) ErrorMessage() string {
	return "Failed to send email."
}

type EmailService struct {
	Controller *EmailController
	View       *EmailView
}

func NewEmailService(controller *EmailController, view *EmailView) *EmailService {
	return &EmailService{
		Controller: controller,
		View:       view,
	}
}

func (s *EmailService) SendEmail(email Email) error {
	auth := smtp.PlainAuth("", s.Controller.From, s.Controller.Password, s.Controller.Host)

	to := []string{email.To}
	msg := []byte(
		"To: " + email.To + "\r\n" +
			"Subject: " + email.Subject + "\r\n" +
			"\r\n" +
			email.Body + "\r\n")

	err := smtp.SendMail(s.Controller.Host+":"+strconv.Itoa(s.Controller.Port), auth, s.Controller.From, to, msg)
	if err != nil {
		return err
	}
	return nil
}


type TrafficViolation struct {
	CarPlate          string  `json:"carPlate"`
	Address           string  `json:"address"`
	Date              string  `json:"date"`
	Violation         string  `json:"violation"`
	CarType           string  `json:"carType"`
	CarColor          string  `json:"carColor"`
	CarBrand          string  `json:"carBrand"`
	VehicleOwnerName  string  `json:"vehicleOwnerName"`
	VehicleOwnerCPF   string  `json:"vehicleOwnerCPF"`
	Speed             float64 `json:"speed"`
	MaxSpeed          int     `json:"maxSpeed"`
	FinePrice         float64 `json:"finePrice"`
	Sex               string  `json:"sex"`
	Age               int  `json:"age"`
}