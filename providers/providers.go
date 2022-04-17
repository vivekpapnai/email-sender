package providers

import (
	"context"
	"emailSender/models"
)

type StorageProvider interface {
	Upload(ctx context.Context, fileName, filePath, contentType string) (string, error)
	GetSharableURL() (string, error)
}

type KafkaProvider interface {
	Publish(topic models.Topic, message []byte)
}
