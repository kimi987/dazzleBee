package main


import (
	"github.com/sirupsen/logrus"
	"os"
)
var log = logrus.New()


func main( ) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.WithFields(logrus.Fields{
        "animal": "walrus",
        "size":   10,
    }).Info("A group of walrus emerges from the ocean")
}