package usersvc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/api"
)

func (svc *UserSvcImpl) SignUp(c context.Context, req api.SignupRequest) (api.SignupResponse, error) {

	var message string

	// print the email and password
	if req.Email != nil {
		fmt.Println(*req.Email)
	}

	if req.Password != nil {
		fmt.Println(*req.Password)
	}

	message = "Signup successful. Please check your email for the verification link."

	var signupResponse api.SignupResponse
	signupResponse.Message = &message

	var msg Message
	msg.From = "anujpflash@gmail.com"
	msg.To = []string{string(*req.Email)}
	msg.Subject = "Signup successful"
	msg.Body = "Please check your email for the verification link."
	msg.Type = "text"

	exchange := "" // Use an empty exchange for direct exchange (default)
	routingKey := "mail"

	// Publish the message to the queue
	body, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}
	err = svc.messageBroker.Publish(c, exchange, routingKey, body)
	if err != nil {
		return signupResponse, err
	}

	return signupResponse, nil
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
