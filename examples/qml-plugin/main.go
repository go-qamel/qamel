package main

import (
	"os"

	"github.com/go-qamel/qamel"
)

func main() {
	// Create application
	app := qamel.NewApplication(len(os.Args), os.Args)
	app.SetApplicationDisplayName("Qamel Example")

	view := qamel.NewViewer()
	engine := view.Engine()
	//$QT_DIR
	engine.AddImportPath("/Users/FlyingtimeICE/QT/Qt5.14.0/Examples/Qt-5.14.0/qml/qmlextensionplugins/imports")

	// Create a QML viewer
	view.SetSource("qrc:/res/main.qml")
	view.SetResizeMode(qamel.SizeRootObjectToView)
	// view.SetHeight(300)
	// view.SetWidth(400)
	view.Show()

	// Exec app
	app.Exec()
}
