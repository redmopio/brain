package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minskylab/brain/channels"
)

func main() {
	client := channels.NewWhatsAppConnector("examplestore.db").Connect()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
