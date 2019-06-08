package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/sylank/lavender-commons-go/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/sylank/lavender-commons-go/properties"

	"log"
	"net/smtp"
)

// SendMailRequest ..
type SendMailRequest struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

// SendConfig ...
type SendConfig struct {
	From     string
	To       string
	Password string
	Message  string
	Subject  string
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	emailSecrets, err := properties.ReadEmailSecrets("config/secrets.json")
	if err != nil {
		log.Fatal("Failed to read email secrets")
	}
	mainTemplate := string(utils.ReadBytesFromFile("config/main_template.template"))

	for _, event := range sqsEvent.Records {
		var request SendMailRequest
		err := json.Unmarshal([]byte(event.Body), &request)
		if err != nil {
			log.Fatal("Failed to unmarshal request")
		}

		mailText := generateMessage(mainTemplate, request.Body)

		sendConfig := &SendConfig{
			From:     emailSecrets.FromAddress,
			To:       request.ToAddress,
			Message:  mailText,
			Subject:  request.Subject,
			Password: emailSecrets.Password,
		}

		sendingError := sendMail(mailText, sendConfig)
		if sendingError != nil {
			log.Println("Failed to send email")
		}
	}

	return nil
}

func generateMessage(template string, body string) string {
	var tmpText = template
	r := strings.NewReplacer("<!-- body -->", body)

	return r.Replace(tmpText)
}

func sendMail(message string, sendConfig *SendConfig) error {
	msg := "From: Levendula Apartman <" + sendConfig.From + ">\n" +
		"To:" + sendConfig.To + "\n" +
		"Subject: " + sendConfig.Subject + "\n\n" +
		sendConfig.Message

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", sendConfig.From, sendConfig.Password, "smtp.gmail.com"),
		sendConfig.From, []string{sendConfig.To}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
