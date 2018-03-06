package main

import (
	"io"
	"net"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	signals := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for _ = range signals {
			log.Info("\nReceived an interrupt, stopping...")
			stop <- true
		}
	}()

	incoming, err := net.Listen("tcp", opts.From)
	if err != nil {
		log.Fatalf("could not start server on %v: %v", opts.From, err)
	}
	log.Infof("server running on %v\n", opts.From)

	client, err := incoming.Accept()
	if err != nil {
		log.Fatal("could not accept client connection", err)
	}
	defer client.Close()
	log.Infof("client '%v' connected!\n", client.RemoteAddr())

	target, err := net.Dial("tcp", opts.To)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	defer target.Close()
	log.Infof("connection to server %v established!\n", target.RemoteAddr())

	go func() { io.Copy(target, client) }()
	go func() { io.Copy(client, target) }()

	<-stop
}
