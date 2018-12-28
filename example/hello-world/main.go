package main

import (
	"os"

	"github.com/RadhiFadlillah/qamel"
)

func main() {
	// Create application
	app := qamel.NewApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("Qamel Example")

	// Create a QML viewer
	view := qamel.NewViewerWithSource("qrc:/res/main.qml")
	view.SetResizeMode(qamel.SizeRootObjectToView)
	view.SetHeight(300)
	view.SetWidth(400)
	view.Show()

	// Exec app
	app.Exec()
}
