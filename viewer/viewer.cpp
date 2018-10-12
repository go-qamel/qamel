#include "viewer.h"
#include <QQuickView>
#include <QString>
#include <QUrl>
#include <QWindow>
#include <QIcon>
#include <QQmlEngine>

void* Viewer_NewViewer() {
    return new QQuickView();
}

void* Viewer_NewViewerWithSource(char* source) {
    QUrl url = QUrl(QString(source));
    return new QQuickView(url);
}

void Viewer_SetSource(void* ptr, char* url) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setSource(QUrl(QString(url)));
}

void Viewer_SetResizeMode(void* ptr, int resizeMode) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setResizeMode(QQuickView::ResizeMode(resizeMode));
}

void Viewer_SetFlags(void* ptr, int flags) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setFlags(Qt::WindowFlags(flags));
}

void Viewer_SetHeight(void* ptr, int height) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setHeight(height);
}

void Viewer_SetWidth(void* ptr, int width) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setWidth(width);
}

void Viewer_SetMaximumHeight(void* ptr, int height) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setMaximumHeight(height);
}

void Viewer_SetMaximumWidth(void* ptr, int width) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setMaximumWidth(width);
}

void Viewer_SetMinimumHeight(void* ptr, int height) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setMinimumHeight(height);
}

void Viewer_SetMinimumWidth(void* ptr, int width) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setMinimumWidth(width);
}

void Viewer_SetOpacity(void* ptr, double opacity) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setOpacity(opacity);
}

void Viewer_SetTitle(void* ptr, char* title) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setTitle(QString(title));
}

void Viewer_SetVisible(void* ptr, bool visible) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setVisible(visible);
}

void Viewer_SetPosition(void* ptr, int x, int y) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setPosition(x, y);
}

void Viewer_SetIcon(void* ptr, char* fileName) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    QIcon icon = QIcon(QString(fileName));
    view->setIcon(icon);
}

void Viewer_Show(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->show();
}

void Viewer_ShowMaximized(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->showMaximized();
}

void Viewer_ShowMinimized(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->showMinimized();
}

void Viewer_ShowFullScreen(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->showFullScreen();
}

void Viewer_ShowNormal(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->showNormal();
}

void Viewer_SetWindowStates(void* ptr, int state) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->setWindowStates(Qt::WindowStates(state));
}

void Viewer_ClearComponentCache(void* ptr) {
    QQuickView *view = static_cast<QQuickView*>(ptr);
    view->engine()->clearComponentCache();
}
