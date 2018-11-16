// +build ignore
package main

import (
	"os"

	"github.com/RadhiFadlillah/qamel"
)

func main() {
	app := qamel.NewApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("HEEYY")

	view := qamel.NewViewerWithSource("main.qml")
	view.SetResizeMode(qamel.SizeRootObjectToView)
	view.SetTitle("It Works")
	view.SetHeight(600)
	view.SetWidth(800)
	view.ShowMaximized()
	app.Exec()
}
