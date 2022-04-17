package kafkaProvider

import (
	"context"
	"emailSender/models"
	"emailSender/providers"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

const (
	csvFileWriterBatchSize = 1
	zipWriterBatchSize     = 1
)

type Logger interface {
	Printf(string, ...interface{})
}

type KafkaProvider struct {
	csvFileWriter *kafka.Writer
	zipFileWriter *kafka.Writer
}

func NewKafkaProvider() providers.KafkaProvider {
	kafkaHost := "localhost:9092"
	l := log.New(os.Stdout, "kafka writer: ", 0)

	// chatWriter is a kafka writer for chat messages
	csvFileWriter := &kafka.Writer{
		Addr:      kafka.TCP(kafkaHost),
		Topic:     string(models.TopicCSVFileUpload),
		Balancer:  &kafka.RoundRobin{},
		BatchSize: csvFileWriterBatchSize,
		Logger:    l,
	}

	zipFileWriter := &kafka.Writer{
		Addr:      kafka.TCP(kafkaHost),
		Topic:     string(models.TopicZipFileUpload),
		Balancer:  &kafka.RoundRobin{},
		BatchSize: zipWriterBatchSize,
		Logger:    l,
	}

	return &KafkaProvider{
		csvFileWriter: csvFileWriter,
		zipFileWriter: zipFileWriter,
	}
}

func (k *KafkaProvider) Publish(topic models.Topic, message []byte) {
	switch topic {
	case models.TopicCSVFileUpload:
		err := k.csvFileWriter.WriteMessages(context.Background(),
			kafka.Message{
				Value: message,
			},
		)
		if err != nil {
			logrus.Errorf("Publish: failed to write csv file upload: %v", err)
		}
	case models.TopicZipFileUpload:
		err := k.zipFileWriter.WriteMessages(context.Background(),
			kafka.Message{
				Value: message,
			},
		)
		if err != nil {
			logrus.Errorf("Publish: failed to write csv file upload: %v", err)
		}
	default:
		logrus.Warn("Trying to publish on wrong topic")
	}
}

func (k *KafkaProvider) Close() {
	if err := k.zipFileWriter.Close(); err != nil {
		logrus.Errorf("error closing kafka chat connection")
	}
}
