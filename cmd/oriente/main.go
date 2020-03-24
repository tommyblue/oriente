package main

import (
	"math/rand"
	"time"

	"github.com/tommyblue/oriente/api"
	"github.com/tommyblue/oriente/oriente"
)

func main() {
	rand.Seed(time.Now().Unix())
	oriente.Initialize()
	api.Run()
}
