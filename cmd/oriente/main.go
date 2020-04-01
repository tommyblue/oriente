package main

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tommyblue/oriente/api"
	"github.com/tommyblue/oriente/oriente"
	"github.com/tommyblue/oriente/store"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	rand.Seed(time.Now().Unix())
	db := store.Init("./db.sql")
	defer db.Close()
	o := oriente.Initialize(db)
	api.Run(o)
}
