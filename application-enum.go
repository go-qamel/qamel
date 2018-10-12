package qamel

// FontWeight is the weight for thickness of the font.
// Qt uses a weighting scale from 0 to 99 similar to,
// but not the same as, the scales used in Windows or CSS.
// A weight of 0 will be thin, whilst 99 will be extremely black.
type FontWeight int

const (
	// Thin is barely visible, usually called hairline
	Thin FontWeight = 0
	// ExtraLight is same with Ultra Light in CSS
	ExtraLight = 12
	// Light is quite visible and almost normal
	Light = 25
	// Normal is the normal font weight
	Normal = 50
	// Medium is a bit heavier than normal
	Medium = 57
	// DemiBold is similar with Semi Bold in CSS
	DemiBold = 63
	// Bold is the bold font, usually used for emphasis
	Bold = 75
	// ExtraBold is usually called Heavy, useful for title
	ExtraBold = 81
	// Black is the thickest font weight, usually used for poster
	Black = 87
)

// Attribute describes attributes that change the behavior of application-wide features.
type Attribute int

const (
	// DontShowIconsInMenus make actions with the Icon property
	// won't be shown in any menus unless specifically set by
	// the QAction::iconVisibleInMenu property.
	// Menus that are currently open or menus already created in
	// the native macOS menubar may not pick up a change in this
	// attribute. Changes in the QAction::iconVisibleInMenu property
	// will always be picked up.
	DontShowIconsInMenus Attribute = 2

	// DontShowShortcutsInContextMenus make actions with the Shortcut
	// property won't be shown in any shortcut menus unless specifically
	// set by the QAction::shortcutVisibleInContextMenu property.
	// This value was added in Qt 5.10.
	DontShowShortcutsInContextMenus = 28

	// NativeWindows ensures that widgets have native windows.
	NativeWindows = 3

	// DontCreateNativeWidgetSiblings ensures that siblings of native
	// widgets stay non-native unless specifically set by the
	// Qt::WA_NativeWindow attribute.
	DontCreateNativeWidgetSiblings = 4

	// PluginApplication indicates that Qt is used to author a plugin.
	// Depending on the operating system, it suppresses specific
	// initializations that do not necessarily make sense in the plugin case.
	// For example on macOS, this includes avoiding loading our nib for the
	// main menu and not taking possession of the native menu bar. Setting
	// this attribute to true will also set the AA_DontUseNativeMenuBar
	// attribute to true. It also disables native event filters. This
	// attribute must be set before QGuiApplication constructed.
	// This value was added in Qt 5.7.
	PluginApplication = 5

	// DontUseNativeMenuBar make all menubars created while this attribute
	// is set to true won't be used as a native menubar (e.g, the menubar
	// at the top of the main screen on macOS).
	DontUseNativeMenuBar = 6

	// MacDontSwapCtrlAndMeta prevents Qt swaps the Control and Meta (Command)
	// keys in MacOS. Whenever Control is pressed, Qt sends Meta, and whenever
	// Meta is pressed Control is sent. When this attribute is true, Qt will
	// not do the flip. QKeySequence::StandardKey will also flip accordingly
	// (i.e., QKeySequence::Copy will be Command+C on the keyboard regardless
	// of the value set, though what is output for QKeySequence::toString()
	// will be different).
	MacDontSwapCtrlAndMeta = 7

	// Use96Dpi assumes the screen has a resolution of 96 DPI rather than
	// using the OS-provided resolution. This will cause font rendering to be
	// consistent in pixels-per-point across devices rather than defining
	// 1 point as 1/72 inch.
	Use96Dpi = 8

	// SynthesizeTouchForUnhandledMouseEvents make all mouse events that are
	// not accepted by the application will be translated to touch events instead.
	SynthesizeTouchForUnhandledMouseEvents = 11

	// SynthesizeMouseForUnhandledTouchEvents make all touch events that are
	// not accepted by the application will be translated to left button mouse
	// events instead. This attribute is enabled by default.
	SynthesizeMouseForUnhandledTouchEvents = 12

	// UseHighDpiPixmaps make QIcon::pixmap() generate high-dpi pixmaps that
	// can be larger than the requested size. Such pixmaps will have
	// devicePixelRatio() set to a value higher than 1. After setting this
	// attribute, application code that uses pixmap sizes in layout geometry
	// calculations should typically divide by devicePixelRatio() to get
	// device-independent layout geometry.
	UseHighDpiPixmaps = 13

	// ForceRasterWidgets make top-level widgets use pure raster surfaces,
	// and do not support non-native GL-based child widgets.
	ForceRasterWidgets = 14

	// UseDesktopOpenGL forces the usage of desktop OpenGL (for example,
	// opengl32.dll or libGL.so) on platforms that use dynamic loading of
	// the OpenGL implementation. This attribute must be set before
	// QGuiApplication is constructed. This value was added in Qt 5.3.
	UseDesktopOpenGL = 15

	// UseOpenGLES forces the usage of OpenGL ES 2.0 or higher on platforms
	// that use dynamic loading of the OpenGL implementation. This attribute
	// must be set before QGuiApplication is constructed.
	// This value was added in Qt 5.3.
	UseOpenGLES = 16

	// UseSoftwareOpenGL forces the usage of a software based OpenGL
	// implementation on platforms that use dynamic loading of the OpenGL
	// implementation. This will typically be a patched build of Mesa llvmpipe,
	// providing OpenGL 2.1. The value may have no effect if no such OpenGL
	// implementation is available. The default name of this library is
	// opengl32sw.dll and can be overridden by setting the environment variable
	// QT_OPENGL_DLL. See the platform-specific pages, for instance Qt for
	// Windows, for more information. This attribute must be set before
	// QGuiApplication is constructed. This value was added in Qt 5.4.
	UseSoftwareOpenGL = 17

	// ShareOpenGLContexts enables resource sharing between the OpenGL contexts
	// used by classes like QOpenGLWidget and QQuickWidget. This allows sharing
	// OpenGL resources, like textures, between QOpenGLWidget instances that
	// belong to different top-level windows. This attribute must be set
	// before QGuiApplication is constructed. This value was added in Qt 5.4.
	ShareOpenGLContexts = 18

	// SetPalette indicates whether a palette was explicitly set on the
	// QGuiApplication. This value was added in Qt 5.5.
	SetPalette = 19

	// EnableHighDpiScaling enables high-DPI scaling in Qt on supported platforms
	// (see also High DPI Displays). Supported platforms are X11, Windows and
	// Android. Enabling makes Qt scale the main (device independent) coordinate
	// system according to display scale factors provided by the operating system.
	// This corresponds to setting the QT_AUTO_SCREEN​_SCALE_FACTOR environment
	// variable to 1. This attribute must be set before QGuiApplication is
	// constructed. This value was added in Qt 5.6.
	EnableHighDpiScaling = 20

	// DisableHighDpiScaling disables high-DPI scaling in Qt, exposing window
	// system coordinates. Note that the window system may do its own scaling,
	// so this does not guarantee that QPaintDevice::devicePixelRatio() will be
	// equal to 1. In addition, scale factors set by QT_SCALE_FACTOR will not be
	// affected. This corresponds to setting the QT_AUTO_SCREEN​_SCALE_FACTOR
	// environment variable to 0. This attribute must be set before QGuiApplication
	// is constructed. This value was added in Qt 5.6.
	DisableHighDpiScaling = 21

	// UseStyleSheetPropagationInWidgetStyles enables Qt Style Sheet in regular QWidget.
	// By default, Qt Style Sheets disable regular QWidget palette and font propagation.
	// When this flag is enabled, font and palette changes propagate as though the
	// user had manually called the corresponding QWidget methods. See The Style
	// Sheet Syntax - Inheritance for more details. This value was added in Qt 5.7.
	UseStyleSheetPropagationInWidgetStyles = 22

	// DontUseNativeDialogs makes all dialogs created while this attribute is
	// set to true won't use the native dialogs provided by the platform.
	// This value was added in Qt 5.7.
	DontUseNativeDialogs = 23

	// SynthesizeMouseForUnhandledTabletEvents makes all tablet events that are not
	// accepted by the application will be translated to mouse events instead.
	// This attribute is enabled by default. This value was added in Qt 5.7.
	SynthesizeMouseForUnhandledTabletEvents = 24

	// CompressHighFrequencyEvents enables compression of certain frequent events.
	// On the X11 windowing system, the default value is true, which means that
	// QEvent::MouseMove, QEvent::TouchUpdate, and changes in window size and position
	// will be combined whenever they occur more frequently than the application
	// handles them, so that they don't accumulate and overwhelm the application later.
	// On other platforms, the default is false. (In the future, the compression
	// feature may be implemented across platforms.) You can test the attribute to
	// see whether compression is enabled. If your application needs to handle all
	// events with no compression, you can unset this attribute. Notice that input
	// events from tablet devices will not be compressed. See AA_CompressTabletEvents
	// if you want these to be compressed as well. This value was added in Qt 5.7.
	CompressHighFrequencyEvents = 25

	// CompressTabletEvents enables compression of input events from tablet devices.
	// Notice that AA_CompressHighFrequencyEvents must be true for events compression
	// to be enabled, and that this flag extends the former to tablet events.
	// Its default value is false. This value was added in Qt 5.10.
	CompressTabletEvents = 29

	// DontCheckOpenGLContextThreadAffinity makes a context that created using
	// QOpenGLContext does not check that the QObject thread affinity of the
	// QOpenGLContext object is the same thread calling makeCurrent().
	// This value was added in Qt 5.8.
	DontCheckOpenGLContextThreadAffinity = 26

	// DisableShaderDiskCache disables caching of shader program binaries on disk.
	// By default Qt Quick, QPainter's OpenGL backend, and any application using
	// QOpenGLShaderProgram with one of its addCacheableShaderFromSource overloads
	// will employ a disk-based program binary cache in either the shared or
	// per-process cache storage location, on systems that support glProgramBinary().
	// In the unlikely event of this being problematic, set this attribute to
	// disable all disk-based caching of shaders.
	DisableShaderDiskCache = 27

	// DisableWindowContextHelpButton disables the WindowContextHelpButtonHint by
	// default on Qt::Sheet and Qt::Dialog widgets. This hides the ? button on
	// Windows, which only makes sense if you use QWhatsThis functionality.
	// This value was added in Qt 5.10. In Qt 6, WindowContextHelpButtonHint
	// will not be set by default.
	DisableWindowContextHelpButton = 30
)
