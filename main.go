package main

import (
	"emailSender/server"
	"fmt"
)

func main() {
	srv := server.SrvInit()
	fmt.Println("connected")
	go srv.Subscribe()
	srv.Start(":8081")
}
