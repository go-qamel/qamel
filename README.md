Qamel
-----

Qamel is a simple QML binding for Go, heavily inspired by [`therecipe/qt`](https://github.com/therecipe/qt). This package only binds Qt's classes that used for creating a simple QML app, i.e. `QApplication`, `QQuickView` and `QQMLApplicationEngine`. It's still work-in progress, however it should be stable enough to use in production (as in I'm already using it in prod without problem, your situations may vary).

### Features

- Published under MIT License, which means you can use this binding to create proprietary app.
- Since it only binds the small set of Qt's class, the build time is quite fast.
- The binding itself is really simple and small. I also think I did a good job on commenting my code, so people should be able to fork it easily.

### Limitation

- I've only tested this in Linux, so I can't vouch for other OS. It should works though, since the code itself is really simple.
- When declaring custom QML object, this binding only supports basic data type, i.e. `int`, `int32`, `int64`, `float32`, `float64`, `bool` and `string`. For other data type like struct, array or map, you have to use `string` type and pass it as JSON value.
- Thanks to Go and Qt, in theory, the app built using this binding can be cross compiled from and to Windows, Linux and MacOS. However, since I only have Linux PC, I only able to test cross compiling from Linux to Windows.

### Licenses

Qamel is distributed under [MIT license](https://choosealicense.com/licenses/mit/), which means you can use and modify it however you want. However, if you make an enhancement for it, if possible, please send a pull request.