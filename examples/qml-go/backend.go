package main

import (
	"math/rand"
	"time"

	"github.com/go-qamel/qamel"
)

// BackEnd is the bridge for communicating between QML and Go
type BackEnd struct {
	qamel.QmlObject
	_ func() int `slot:"getRandomNumber"`
}

func (b *BackEnd) getRandomNumber() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(9999)
}
