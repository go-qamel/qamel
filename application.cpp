#include "application.h"
#include <QCoreApplication>
#include <QGuiApplication>
#include <QFont>
#include <QString>
#include <QIcon>

void* App_NewApplication(int argc, char* argv[]) {
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    return new QGuiApplication(argc, argv);
}

void App_SetAttribute(int attribute, bool on) {
    QCoreApplication::setAttribute(Qt::ApplicationAttribute(attribute), on);
}

void App_SetFont(void *ptr, char *family, int pointSize, int weight, bool italic) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    QFont font = QFont(QString(family), pointSize, weight, italic);
    app->setFont(font);
}

void App_SetQuitOnLastWindowClosed(void* ptr, bool quit) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setQuitOnLastWindowClosed(quit);
}

void App_SetApplicationDisplayName(void* ptr, char* name) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setApplicationDisplayName(QString(name));
}

void App_SetWindowIcon(void* ptr, char* fileName) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    QIcon icon = QIcon(QString(fileName));
    app->setWindowIcon(icon);
}

void App_SetApplicationName(void* ptr, char* name) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setApplicationName(QString(name));
}

void App_SetApplicationVersion(void* ptr, char* version) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setApplicationVersion(QString(version));
}

void App_SetOrganizationName(void* ptr, char* name) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setOrganizationName(QString(name));
}

void App_SetOrganizationDomain(void* ptr, char* domain) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    app->setOrganizationDomain(QString(domain));
}

int App_Exec(void* ptr) {
    QGuiApplication *app = static_cast<QGuiApplication*>(ptr);
    return app->exec();
}