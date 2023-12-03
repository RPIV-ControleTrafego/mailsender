package main
//+build 
import (
	"net/smtp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	
	"github.com/arsmn/fiber-swagger/v2"
	
)


// @title Email Service API
// @version 1.0
// @description API for sending emails using Fiber
// @host localhost:3000
// @BasePath /
func main() {
	// Criação e inicialização do Service Locator
	serviceLocator := NewEmailServiceLocator()
	serviceLocator.Initialize(
		"andremendes0113@gmail.com",
		"0294BE2A8D274556D9584EE59D90DEFC0AB6",
		"smtp.elasticemail.com",
		2525,
	)

	// Configuração do aplicativo
	app := setupApp(serviceLocator)

	// Inicialização do aplicativo
	app.Listen(":3000")
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
