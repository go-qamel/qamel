#include "engine.h"
#include <QQmlApplicationEngine>
#include <QString>
#include <QUrl>

void* Engine_NewEngine() {
    return new QQmlApplicationEngine();
}

void Engine_Load(void* ptr, char* url) {
    QQmlApplicationEngine *engine = static_cast<QQmlApplicationEngine*>(ptr);
    engine->load(QUrl(QString(url)));
}

void Engine_ClearComponentCache(void* ptr) {
    QQmlApplicationEngine *engine = static_cast<QQmlApplicationEngine*>(ptr);
    engine->clearComponentCache();
}
