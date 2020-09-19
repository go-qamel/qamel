import QtQuick 2.12
import TimeExample 1.0 // import types from the plugin

Clock { // this class is defined in QML (imports/TimeExample/Clock.qml)
    Time { // this class is defined in C++ (plugin.cpp)
        id: time
    }

    hours: time.hour
    minutes: time.minute
}
