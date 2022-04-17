package utils

import (
	"emailSender/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

func EncodeJSON200Body(resp http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
}

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
}

func DownLoadImageFileFromURL(uploadImage models.UploadImage) (string, string, error) {
	response, err := http.Get(uploadImage.URL)
	if err != nil {
		return "", "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != 200 {
		return "", "", errors.New("received non 200 response code")
	}

	//imageName := fmt.Sprintf("image-%v", time.Now().Unix())
	fileName := fmt.Sprintf("%v-%v.jpg", uploadImage.Name, time.Now().Unix())
	file, err := os.Create(fileName)
	if err != nil {
		return "", "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	contentType := response.Header.Get("Content-Type")
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", contentType, err
	}

	return fileName, contentType, err
}
