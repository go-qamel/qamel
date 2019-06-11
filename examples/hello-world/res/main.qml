import QtQuick 2.12

Rectangle {
    color: "cyan"

    Text {
        anchors.fill: parent
        text: "Hello World"
        font.pixelSize: 32
        font.weight: Font.Bold
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
    }
}