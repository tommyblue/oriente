package main

import (
	"math/rand"
	"time"

	"github.com/tommyblue/oriente/oriente"
)

func main() {
	rand.Seed(time.Now().Unix())
	oriente.NewGame(2)
}
