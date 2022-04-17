package server

import (
	"context"
	"emailSender/models"
	"emailSender/utils"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

func (srv *Server) greet(resp http.ResponseWriter, _ *http.Request) {
	utils.EncodeJSON200Body(resp, map[string]interface{}{
		"message": "welcome to image downloading service",
	})
}

type Logger interface {
	Printf(string, ...interface{})
}

func (srv *Server) Subscribe() {
	kafkaHost := "localhost:9092"

	l := log.New(os.Stdout, "kafka reader: ", 0)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaHost},
		Topic:   string(models.TopicZipFileUpload),
		GroupID: "zip-consumer-group",
		Logger:  l,
	})

	type sendZipVFile struct {
		URL   string `json:"data"`
		Email string `json:"email"`
	}

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		logrus.Infof("zipFile: time %v read message at offset %d: %s = %s\n", time.Now(), m.Offset, string(m.Key), string(m.Value))

		var message sendZipVFile
		if err = json.Unmarshal(m.Value, &message); err != nil {
			logrus.Errorf("Failed to Unmarshal message from topic: %q error: %v", m.Topic, err)
			continue
		}

		go sendZip(message.URL, message.Email)

		fmt.Printf("data in message is: %v\n", message)

	}
}

func sendZip(url string, email string) {
	from := mail.NewEmail("vivek papnai", "papnaivivek@gmail.com")
	subject := "Download Your Zip File HERE"
	to := mail.NewEmail("User", email)
	plainTextContent := "send Email"
	htmlContent := fmt.Sprintf("<strong>Hey! Your zip file is ready please use the link given below to download it<br>url:%v", url)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SEND_GRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		//fmt.Println(response.Body)
		//fmt.Println(response.Headers)
	}
}
