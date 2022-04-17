package models

type UploadImage struct {
	Name string
	URL  string
}

type Topic string
type KafkaHeaders string

const (
	TopicCSVFileUpload Topic = "csv_file"
	TopicZipFileUpload Topic = "zip_file"

	KafkaHeadersWSConnectionUnixNano KafkaHeaders = "connection_time_unix_nano"
)
