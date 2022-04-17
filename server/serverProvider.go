package server

import (
	"emailSender/providers"
	"emailSender/providers/kafkaProvider"
	"emailSender/providers/storageProvider"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	defaultServerRequestTimeoutMinutes      = 2
	defaultServerReadHeaderTimeoutSeconds   = 30
	defaultServerRequestWriteTimeoutMinutes = 30
)

type Server struct {
	StorageProvider providers.StorageProvider
	httpServer      *http.Server
	KafkaProvider   providers.KafkaProvider
}

func SrvInit() *Server {
	// storage provider is the storage used to upload files
	sp := storageProvider.NewStorageProvider()

	// kafkaProvider is a PubSub and queue service for the notification, chats etc.
	kp := kafkaProvider.NewKafkaProvider()

	return &Server{StorageProvider: sp,
		KafkaProvider: kp}
}

func (srv *Server) Start(addr string) {
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.InjectRoutes(),
		ReadTimeout:       defaultServerRequestTimeoutMinutes * time.Minute,
		ReadHeaderTimeout: defaultServerReadHeaderTimeoutSeconds * time.Second,
		WriteTimeout:      defaultServerRequestWriteTimeoutMinutes * time.Minute,
	}
	srv.httpServer = httpSrv

	logrus.Info("Server running at PORT ", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
		return
	}
}
