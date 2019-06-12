package main

import (
	"github.com/RadhiFadlillah/qamel/internal/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.QamelCmd().Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
