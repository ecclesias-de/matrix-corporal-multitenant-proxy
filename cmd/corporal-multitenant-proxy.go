package main

import (
	"flag"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/container"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "", "configuration file to use")
	flag.Parse()

	config, err := config.Load(*configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	services := container.Bootstrap(config)

	for _, service := range services {
		err := service.Start()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	shutdown := make(chan bool)

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel

		for _, service := range services {
			err := service.Stop()
			if err != nil {
				logrus.Warn(err)
			}
		}

		shutdown <- true
	}()

	<-shutdown
}
