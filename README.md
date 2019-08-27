Qamel
-----

[![GoDoc](https://godoc.org/github.com/RadhiFadlillah/qamel?status.png)](https://godoc.org/github.com/RadhiFadlillah/qamel)
[![Docker](https://img.shields.io/badge/docker-qamel-blue.svg)](https://hub.docker.com/r/radhifadlillah/qamel)
[![Donate](https://img.shields.io/badge/donate-PayPal-green.svg)](https://www.paypal.me/RadhiFadlillah)
[![Donate](https://img.shields.io/badge/donate-Ko--fi-brightgreen)](https://ko-fi.com/radhifadlillah)

Qamel is a simple QML binding for Go, heavily inspired by [`therecipe/qt`](https://github.com/therecipe/qt). This package only binds Qt's classes that used for creating a simple QML app, i.e. `QApplication`, `QQuickView` and `QQMLApplicationEngine`. It's still work-in progress, however it should be stable enough to use in production (as in I'm already using it in prod without problem, your situations may vary).

### Features

- Published under MIT License, which means you can use this binding for whatever you want.
- Since it only binds the small set of Qt's class, the build time is quite fast.
- It's available as Docker image, which means you can create QML app without installing Qt in your PC. Go is still needed though.
- The binding itself is really simple and small. I also think I did a good job on commenting my code, so people should be able to fork it easily.
- It supports [live reload](https://godoc.org/github.com/RadhiFadlillah/qamel#Viewer.WatchResourceDir) which is really useful while working on GUI.

### Limitation

- I've only tested this in Linux and Windows, so I'm not sure about Mac OS. It should works though, since the code itself is really simple.
- When declaring custom QML object, this binding only [supports](https://github.com/RadhiFadlillah/qamel/wiki/QmlObject-Documentation) basic data type, i.e. `int`, `int32`, `int64`, `float32`, `float64`, `bool` and `string`. For other data type like struct, array or map, you have to use `string` type and pass it as JSON value.
- Thanks to Go and Qt, in theory, the app built using this binding can be cross compiled from and to Windows, Linux and MacOS. However, since I only have Linux and Windows PC, I only able to test cross compiling between Linux and Windows.

### Development Status

I've created this binding for my job, so it's actively maintained. However, since I created it for the sake of the job, if the issues are not critical and doesn't affect my job or workflow, it might take a long time before I work on it. Therefore, all PRs and contributors will always be welcomed.

### Resources

All documentation for this binding is available in [wiki](https://github.com/RadhiFadlillah/qamel/wiki) and [GoDoc](https://godoc.org/github.com/RadhiFadlillah/qamel). There are some important sections in wiki that I recommend you to check before you start developing your QML app :

- [Frequently Asked Questions](https://github.com/RadhiFadlillah/qamel/wiki/Frequently-Asked-Questions-(FAQ)), especially the [comparison](https://github.com/RadhiFadlillah/qamel/wiki/Frequently-Asked-Questions-(FAQ)#how-does-it-compare-to-therecipeqt-) with other binding.
- [Getting Started](https://github.com/RadhiFadlillah/qamel/wiki/Getting-Started)
- [CLI Usage](https://github.com/RadhiFadlillah/qamel/wiki/CLI-Usage)
- [Building App](https://github.com/RadhiFadlillah/qamel/wiki/Building-Application)

You might also want to check Qt's official documentation about QML :

- [Qt QML](http://doc.qt.io/qt-5/qtqml-index.html)
- [All QML Types](http://doc.qt.io/qt-5/qmltypes.html)
- [QML Reference](http://doc.qt.io/qt-5/qmlreference.html)
- [QML Code Conventions](https://doc.qt.io/qt-5/qml-codingconventions.html)

For demo, you can check out [Qamel-HN](https://github.com/RadhiFadlillah/qamel-hn), a HackerNews reader built with QML and Go.

### Licenses

Qamel is distributed under [MIT license](https://choosealicense.com/licenses/mit/), which means you can use and modify it however you want. However, if you make an enhancement for it, if possible, please send a pull request. If you like this project, please consider donating to me either via [PayPal](https://www.paypal.me/RadhiFadlillah) or [Ko-Fi](https://ko-fi.com/radhifadlillah).
