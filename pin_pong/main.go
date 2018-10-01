package main

import (
	"time"
	"github.com/Sirupsen/logrus"
)

var message = ""
var reply = ""


func pinger(pinger chan string) {
	for {
		message = <-pinger
		time.Sleep(time.Second)
		if message == "ping" {
			reply = "pong"
			pinger <- reply
		} else {
			reply = "ping"
			pinger <- reply
		}
		logrus.WithFields(logrus.Fields{
			"uid": 1,
			"got": message,
			"sent": reply,
		}).Info("New Round")

	}
}


func ponger(pinger chan string) {
	for {
		message = <-pinger
		time.Sleep(time.Second)
		if message == "ping" {
			reply = "pong"
			pinger <- reply
		} else {
			reply = "ping"
			pinger <- reply
		}
		logrus.WithFields(logrus.Fields{
			"uid": 2,
			"got": message,
			"sent": reply,
		}).Info("New Round")
	}

}

func main() {
	ping := make(chan string, 1)
	ping <- "pong"
	go pinger(ping)
	go ponger(ping)

	// The main goroutine starts the ping/pong by sending into the channel

	for {
		// Block the main thread until an interrupt
		time.Sleep(time.Second)
	}
}