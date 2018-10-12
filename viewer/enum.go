package viewer

// ResizeMode specifies how to resize the view.
type ResizeMode int

const (
	// SizeViewToRootObject makes the view resizes with the root item in the QML.
	SizeViewToRootObject ResizeMode = 0

	// SizeRootObjectToView makes the view will automatically resize
	// the root item to the size of the view.
	SizeRootObjectToView = 1
)

// WindowFlags is used to specify various window-system properties for the widget.
// They are fairly unusual but necessary in a few cases. Some of these flags depend on
// whether the underlying window manager supports them.
type WindowFlags int64

const (
	// Widget is the default type for QWidget. Widgets of this type are child widgets
	// if they have a parent, and independent windows if they have no parent.
	Widget WindowFlags = 0x00000000

	// Window indicates that the widget is a window, usually with a window system frame and
	// a title bar, irrespective of whether the widget has a parent or not. Note that it is
	// not possible to unset this flag if the widget does not have a parent.
	Window = 0x00000001

	// Dialog indicates that the widget is a window that should be decorated as a dialog
	// (i.e., typically no maximize or minimize buttons in the title bar). This is the default
	// type for QDialog. If you want to use it as a modal dialog, it should be launched from
	// another window, or have a parent and used with the QWidget::windowModality property.
	// If you make it modal, the dialog will prevent other top-level windows in the application
	// from getting any input. We refer to a top-level window that has a parent as a secondary window.
	Dialog = 0x00000002 | Window

	// Sheet indicates that the window is a sheet on macOS. Since using a sheet implies window
	// modality, the recommended way is to use QWidget::setWindowModality(), or QDialog::open(), instead.
	Sheet = 0x00000004 | Window

	// Drawer indicates that the widget is a drawer on macOS.
	Drawer = Sheet | Dialog

	// Popup indicates that the widget is a pop-up top-level window, i.e. that it is modal, but has
	// a window system frame appropriate for pop-up menus.
	Popup = 0x00000008 | Window

	// Tool indicates that the widget is a tool window. A tool window is often a small window with
	// a smaller than usual title bar and decoration, typically used for collections of tool buttons.
	// If there is a parent, the tool window will always be kept on top of it. If there isn't a
	// parent, you may consider using WindowStaysOnTopHint as well. If the window system supports it,
	// a tool window can be decorated with a somewhat lighter frame. It can also be combined with
	// FramelessWindowHint. On macOS, tool windows correspond to the NSPanel class of windows. This
	// means that the window lives on a level above normal windows making it impossible to put a normal
	// window on top of it. By default, tool windows will disappear when the application is inactive.
	// This can be controlled by the WA_MacAlwaysShowToolWindow attribute.
	Tool = Popup | Dialog

	// ToolTip indicates that the widget is a tooltip. This is used internally to implement tooltips.
	ToolTip = Popup | Sheet

	// SplashScreen indicates that the window is a splash screen. This is the default type for QSplashScreen.
	SplashScreen = ToolTip | Dialog

	// Desktop indicates that this widget is the desktop. This is the type for QDesktopWidget.
	Desktop = 0x00000010 | Window

	// SubWindow indicates that this widget is a sub-window, such as a QMdiSubWindow widget.
	SubWindow = 0x00000012

	// ForeignWindow indicates that this window object is a handle representing a native platform
	// window created by another process or by manually using native code.
	ForeignWindow = 0x00000020 | Window

	// CoverWindow indicates that the window represents a cover window, which is shown when the
	// application is minimized on some platforms.
	CoverWindow = 0x00000040 | Window

	// MSWindowsFixedSizeDialogHint gives the window a thin dialog border on Windows. This style is
	// traditionally used for fixed-size dialogs.
	MSWindowsFixedSizeDialogHint = 0x00000100

	// MSWindowsOwnDC gives the window its own display context on Windows.
	MSWindowsOwnDC = 0x00000200

	// BypassWindowManagerHint can be used to indicate to the platform plugin that "all"
	// window manager protocols should be disabled. This flag will behave different depending on what
	// operating system the application is running on and what window manager is running. The flag
	// can be used to get a native window with no configuration set.
	BypassWindowManagerHint = 0x00000400

	// X11BypassWindowManagerHint bypass the window manager completely. This results in a borderless
	// window that is not managed at all (i.e., no keyboard input unless you call
	// QWidget::activateWindow() manually).
	X11BypassWindowManagerHint = BypassWindowManagerHint

	// FramelessWindowHint produces a borderless window. The user cannot move or resize a borderless
	// window via the window system. On X11, the result of the flag is dependent on the window manager
	// and its ability to understand Motif and/or NETWM hints. Most existing modern window managers
	// can handle this.
	FramelessWindowHint = 0x00000800

	// NoDropShadowWindowHint disables window drop shadow on supporting platforms.
	NoDropShadowWindowHint = 0x40000000

	// CustomizeWindowHint turns off the default window title hints.
	CustomizeWindowHint = 0x02000000

	// WindowTitleHint gives the window a title bar.
	WindowTitleHint = 0x00001000

	// WindowSystemMenuHint adds a window system menu, and possibly a close button (for example
	// on Mac). If you need to hide or show a close button, it is more portable to use WindowCloseButtonHint.
	WindowSystemMenuHint = 0x00002000

	// WindowMinimizeButtonHint adds a minimize button. On some platforms this implies
	// WindowSystemMenuHint for it to work.
	WindowMinimizeButtonHint = 0x00004000

	// WindowMaximizeButtonHint adds a maximize button. On some platforms this implies
	// WindowSystemMenuHint for it to work.
	WindowMaximizeButtonHint = 0x00008000

	// WindowMinMaxButtonsHint adds a minimize and a maximize button. On some platforms
	// this implies WindowSystemMenuHint for it to work.
	WindowMinMaxButtonsHint = WindowMinimizeButtonHint | WindowMaximizeButtonHint

	// WindowCloseButtonHint adds a close button. On some platforms this implies
	// WindowSystemMenuHint for it to work.
	WindowCloseButtonHint = 0x08000000

	// WindowContextHelpButtonHint adds a context help button to dialogs. On some platforms this
	// implies WindowSystemMenuHint for it to work.
	WindowContextHelpButtonHint = 0x00010000

	// MacWindowToolBarButtonHint on macOS adds a tool bar button (i.e., the oblong button that is
	// on the top right of windows that have toolbars).
	MacWindowToolBarButtonHint = 0x10000000

	// WindowFullscreenButtonHint on macOS adds a fullscreen button.
	WindowFullscreenButtonHint = 0x80000000

	// BypassGraphicsProxyWidget prevents the window and its children from automatically embedding
	// themselves into a QGraphicsProxyWidget if the parent widget is already embedded. You can set
	// this flag if you want your widget to always be a toplevel widget on the desktop, regardless
	// of whether the parent widget is embedded in a scene or not.
	BypassGraphicsProxyWidget = 0x20000000

	// WindowShadeButtonHint adds a shade button in place of the minimize button if the underlying
	// window manager supports it.
	WindowShadeButtonHint = 0x00020000

	// WindowStaysOnTopHint informs the window system that the window should stay on top of all
	// other windows. Note that on some window managers on X11 you also have to pass
	// X11BypassWindowManagerHint for this flag to work correctly.
	WindowStaysOnTopHint = 0x00040000

	// WindowStaysOnBottomHint informs the window system that the window should stay on bottom of
	// all other windows. Note that on X11 this hint will work only in window managers that support
	// _NET_WM_STATE_BELOW atom. If a window always on the bottom has a parent, the parent will
	// also be left on the bottom. This window hint is currently not implemented for macOS.
	WindowStaysOnBottomHint = 0x04000000

	// WindowTransparentForInput informs the window system that this window is used only for
	// output (displaying something) and does not take input. Therefore input events should
	// pass through as if it wasn't there.
	WindowTransparentForInput = 0x00080000

	// WindowOverridesSystemGestures informs the window system that this window implements its
	// own set of gestures and that system level gestures, like for instance three-finger desktop
	// switching, should be disabled.
	WindowOverridesSystemGestures = 0x00100000

	// WindowDoesNotAcceptFocus informs the window system that this window should not
	// receive the input focus.
	WindowDoesNotAcceptFocus = 0x00200000

	// MaximizeUsingFullscreenGeometryHint informs the window system that when maximizing the
	// window it should use as much of the available screen geometry as possible, including areas
	// that may be covered by system UI such as status bars or application launchers. This may
	// result in the window being placed under these system UIs, but does not guarantee it,
	// depending on whether or not the platform supports it. When the flag is enabled the user
	// is responsible for taking QScreen::availableGeometry() into account, so that any UI elements
	// in the application that require user interaction are not covered by system UI.
	MaximizeUsingFullscreenGeometryHint = 0x00400000
)

// WindowStates is used to specify the current state of a top-level window.
type WindowStates int

const (
	// WindowNoState makes the window has no state set (in normal state).
	WindowNoState WindowStates = 0x00000000

	// WindowMinimized makes the window is minimized (i.e. iconified).
	WindowMinimized = 0x00000001

	// WindowMaximized makes the window is maximized with a frame around it.
	WindowMaximized = 0x00000002

	// WindowFullScreen makes the window fills the entire screen without any frame around it.
	WindowFullScreen = 0x00000004
)
