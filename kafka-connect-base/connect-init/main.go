package main

import (
	"log"
	"net/http"
	"os"

	"github.com/connect-init/kcinit"

	"github.com/sirupsen/logrus"
)

// SRVER is set during build
var SRVER string = "<not set>"

func main() {
	logger := logrus.New()
	logger.WithField("StreamReactor", SRVER).Infof("Initiated")

	// Setup config
	err := kcinit.LoadConfig(logger)
	if err != nil {
		os.Exit(1)
	}

	_, err = kcinit.SetupKafkaConnect(logger)
	if err != nil {
		os.Exit(1)
	}
	_, err = kcinit.Service.SetupService(logger)
	if err != nil {
		os.Exit(1)
	}

	err = kcinit.SetConnector(logger)
	if err != nil {
		logger.WithError(err).Error(
			"Unable to create connector.properties.",
		)
		os.Exit(1)
	}

	http.HandleFunc("/api/status", kcinit.StatusListener)
	http.HandleFunc("/api/start", kcinit.StartServiceListener)
	http.HandleFunc("/api/stop", kcinit.StopServiceListener)

	logger.WithFields(logrus.Fields{
		"Component": "HTTP Service",
		"Stage":     "Init",
	}).Info("Listening to 0.0.0.0:8881")
	err = http.ListenAndServe(":8881", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
