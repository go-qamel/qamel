import QtQuick 2.12
import QtQuick.Controls 2.14

Rectangle {
    color: "cyan"

    Text {
        anchors.fill: parent
        text: "Hello World haha"
        font.pixelSize: 32
        font.weight: Font.Bold
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
    }

    Button {
        text: "ceshi"
    }
}