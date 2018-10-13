#include <QQuickView>
#include <QString>
#include <QUrl>
#include <QWindow>
#include <QIcon>
#include <QQmlEngine>
#include "viewer.h"

class QamelView : public QQuickView {
    Q_OBJECT

public:
    QamelView(QWindow *parent = 0) : QQuickView(parent) {}
    QamelView(const QUrl &source, QWindow *parent = nullptr) : QQuickView(source, parent) {}

public slots:
    void reload() {
        engine()->clearComponentCache();
        setSource(source());
    }
};

void* Viewer_NewViewer() {
    return new QamelView();
}

void* Viewer_NewViewerWithSource(char* source) {
    QUrl url = QUrl(QString(source));
    return new QamelView(url);
}

void Viewer_SetSource(void* ptr, char* url) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setSource(QUrl(QString(url)));
}

void Viewer_SetResizeMode(void* ptr, int resizeMode) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setResizeMode(QQuickView::ResizeMode(resizeMode));
}

void Viewer_SetFlags(void* ptr, int flags) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setFlags(Qt::WindowFlags(flags));
}

void Viewer_SetHeight(void* ptr, int height) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setHeight(height);
}

void Viewer_SetWidth(void* ptr, int width) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setWidth(width);
}

void Viewer_SetMaximumHeight(void* ptr, int height) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setMaximumHeight(height);
}

void Viewer_SetMaximumWidth(void* ptr, int width) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setMaximumWidth(width);
}

void Viewer_SetMinimumHeight(void* ptr, int height) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setMinimumHeight(height);
}

void Viewer_SetMinimumWidth(void* ptr, int width) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setMinimumWidth(width);
}

void Viewer_SetOpacity(void* ptr, double opacity) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setOpacity(opacity);
}

void Viewer_SetTitle(void* ptr, char* title) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setTitle(QString(title));
}

void Viewer_SetVisible(void* ptr, bool visible) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setVisible(visible);
}

void Viewer_SetPosition(void* ptr, int x, int y) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setPosition(x, y);
}

void Viewer_SetIcon(void* ptr, char* fileName) {
    QamelView *view = static_cast<QamelView*>(ptr);
    QIcon icon = QIcon(QString(fileName));
    view->setIcon(icon);
}

void Viewer_Show(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->show();
}

void Viewer_ShowMaximized(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->showMaximized();
}

void Viewer_ShowMinimized(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->showMinimized();
}

void Viewer_ShowFullScreen(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->showFullScreen();
}

void Viewer_ShowNormal(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->showNormal();
}

void Viewer_SetWindowStates(void* ptr, int state) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->setWindowStates(Qt::WindowStates(state));
}

void Viewer_ClearComponentCache(void* ptr) {
    QamelView *view = static_cast<QamelView*>(ptr);
    view->engine()->clearComponentCache();
}

void Viewer_Reload(void* ptr) {
    QMetaObject::invokeMethod(static_cast<QamelView*>(ptr), "reload");
}

#include "moc-viewer.h"