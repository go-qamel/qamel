package main

import (
	"time"

	"github.com/RadhiFadlillah/qamel"
)

// BackEnd is the bridge for communicating between QML and Go
type BackEnd struct {
	qamel.QmlObject
	_ func()       `constructor:"init"`
	_ func(string) `signal:"timeChanged"`
}

func (b *BackEnd) init() {
	go func() {
		for {
			now := time.Now().Format("15:04:05")
			b.timeChanged(now)
			time.Sleep(time.Second)
		}
	}()
}
