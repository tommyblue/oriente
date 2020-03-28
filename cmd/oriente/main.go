package main

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tommyblue/oriente/api"
	"github.com/tommyblue/oriente/oriente"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	rand.Seed(time.Now().Unix())
	oriente.Initialize()
	api.Run()
}
