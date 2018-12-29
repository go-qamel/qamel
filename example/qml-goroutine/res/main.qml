import QtQuick 2.12
import BackEnd 1.0

Rectangle {
    color: "white"

    BackEnd {
        id: backEnd
        onTimeChanged: (time) => txt.text = time
    }

    Text {
        id: txt
        anchors.fill: parent
        font.pixelSize: 32
        font.weight: Font.Bold
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
    }
}
