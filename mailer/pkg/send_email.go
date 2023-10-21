package pkg

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/GDGVIT/opengraph-thumbnail-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type MailInstance struct {
	Logger        logger.Logger
	smtp_username string `validate:"required"`
	smtp_password string `validate:"required"`
	smtp_host     string `validate:"required,ip|hostname"`
	smtp_port     int    `validate:"required,numeric,oneof=25 465 587 2525"`
}
type Message struct {
	From         string
	To           []string
	Subject      string
	Body         string
	TemplateName string
	Data         interface{}
	Type         string // template or text or html
}

// @Function: SetCredentials
// @Description: This function is used to send an email
// @Params: username string - username of the sender
// @Params: password string - password of the sender
// @Return: error (error if any)
func (svc *MailInstance) SetCredentials(username string, password string) error {
	if err := validateVar(username, "required"); err != nil {
		return err
	}
	svc.smtp_username = username

	if err := validateVar(password, "required"); err != nil {
		return err
	}
	svc.smtp_password = password
	return nil
}

// @Function: SetTransportDetails
// @Description: This function is used to set the transport details
// @Params: host string - host of the smtp server
// @Params: port int - port of the smtp server
// @Return: error (error if any)
func (svc *MailInstance) SetTransportDetails(host string, port int) error {
	if err := validateVar(host, "required,ip|hostname"); err != nil {
		svc.Logger.Error(errors.Wrap(err, "[SetTransportDetails] [validateVar]"))
		return err
	}
	svc.smtp_host = host
	if err := validateVar(port, "required,numeric,oneof=25 465 587 2525"); err != nil {
		svc.Logger.Error(errors.Wrap(err, "[SetTransportDetails] [validateVar]"))
		return err
	}
	svc.smtp_port = port
	return nil
}

// @Function: SendEmail
// @Description: This function is used to send an email using a template
func (svc *MailInstance) SendEmail(msg Message) error {
	if err := validateStruct(svc); err != nil {
		svc.Logger.Error(errors.Wrap(err, "[SendMailUsingTemplate] [validateStruct]"))
		return errors.New("invalid credentials")
	}

	// Validate recipients
	if len(msg.To) < 1 {
		return errors.New("at least one recipient required")
	}
	for _, recipient := range msg.To {
		if err := validateVar(recipient, "required,email,min=4"); err != nil {
			svc.Logger.Error(errors.Wrap(err, "[SendMailUsingTemplate] [validateVar]"))
			return err
		}
	}

	// Validate subject
	msg.Subject = strings.TrimSpace(msg.Subject)
	if msg.Subject == "" {
		return errors.New("subject cannot be empty")
	}

	// Prepare Dialer
	host := svc.smtp_host
	port := svc.smtp_port
	username := svc.smtp_username
	password := svc.smtp_password
	dialer := gomail.NewDialer(host, port, username, password)

	// Prepare Message
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	svc.Logger.Info(fmt.Sprintf("Sending email to %s", msg.To))

	switch msg.Type {
	case "template":
		// Parse Template
		body, err := getParsedTemplate(msg.TemplateName, msg.Data)
		if err != nil {
			return err
		}
		m.SetBody("text/html", body)
	case "text":
		m.SetBody("text/plain", msg.Body)
	case "html":
		m.SetBody("text/html", msg.Body)
	default:
		return errors.New("invalid type")
	}
	if err := dialer.DialAndSend(m); err != nil {
		svc.Logger.Error(errors.Wrap(err, "[SendMailUsingTemplate] [dialer.DialAndSend]"))
		return err
	}
	return nil
}

// @Function: GetTransportDetails
// @Description: This function is used to get the transport details
// @Return: host string - host of the smtp server
// @Return: port int - port of the smtp server
func (svc *MailInstance) GetTransportDetails() (string, int) {
	return svc.smtp_host, svc.smtp_port
}

// @Function: GetCredentails
// @Description: This function is used to get the credentials
// @Return: username string - username of the sender
// @Return: password string - password of the sender
func (svc *MailInstance) GetCredentials() (string, string) {
	return svc.smtp_username, svc.smtp_password
}

// @Function: getParsedTemplate
// @Description: This function is used to parse a template
// @Params: filename string - name of the template file
// @Params: data interface{} - data to be passed to the template
// @Return: templateString string - parsed template
// @Return: error (error if any)
func getParsedTemplate(filename string, data interface{}) (string, error) {
	// Working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	templateName := strings.Split(filename, "/")

	// Parse template
	path := fmt.Sprintf("%s/%s", dir, filename)
	templateString, err := template.New(templateName[len(templateName)-1]).ParseFiles(path)
	if err != nil {
		return "", err
	}

	// Execute template
	var templateBuffer bytes.Buffer
	err = templateString.Execute(&templateBuffer, data)
	if err != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

// @Function: ValidateStruct
// @Description: This function is used to validate a variable with a struct
// @Params: variable interface{} - variable to be validated
// @Return: errors []ErrorResponse (errors if any)
func validateStruct(variable interface{}) error {
	var errors *multierror.Error
	validate := validator.New()
	err := validate.Struct(variable)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

// @Function: ValidateVar
// @Description: This function is used to validate a variable
// @Params: field interface{} - variable to be validated
// @Params: tag string - tag to be used
// @Return: error (error if any)
func validateVar(field interface{}, tag string) error {
	validate := validator.New()
	err := validate.Var(field, tag)
	return err
}
