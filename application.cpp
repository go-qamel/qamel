#include "application.h"
#include <QGuiApplication>
#include <QFont>
#include <QString>
#include <QIcon>
#include <QByteArray>
#include <QList>

void* App_NewApplication(int argc, char* argv) {
	static int argcs = argc;
	static char** argvs = static_cast<char**>(malloc(argcs * sizeof(char*)));

	QList<QByteArray> aList = QByteArray(argv).split('|');
	for (int i = 0; i < argcs; i++) {
		argvs[i] = (new QByteArray(aList.at(i)))->data();
	}

    // QGuiApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    return new QGuiApplication(argcs, argvs);
}

void App_SetAttribute(long long attribute, bool on) {
    QGuiApplication::setAttribute(Qt::ApplicationAttribute(attribute), on);
}

void App_SetFont(char *family, int pointSize, int weight, bool italic) {
    QFont font = QFont(QString(family), pointSize, weight, italic);
    QGuiApplication::setFont(font);
}

void App_SetQuitOnLastWindowClosed(bool quit) {
    QGuiApplication::setQuitOnLastWindowClosed(quit);
}

void App_SetApplicationDisplayName(char* name) {
	QGuiApplication::setApplicationDisplayName(QString(name));
}

void App_SetWindowIcon(char* fileName) {
    QIcon icon = QIcon(QString(fileName));
    QGuiApplication::setWindowIcon(icon);
}

void App_SetApplicationName(char* name) {
    QGuiApplication::setApplicationName(QString(name));
}

void App_SetApplicationVersion(char* version) {
    QGuiApplication::setApplicationVersion(QString(version));
}

void App_SetOrganizationName(char* name) {
    QGuiApplication::setOrganizationName(QString(name));
}

void App_SetOrganizationDomain(char* domain) {
    QGuiApplication::setOrganizationDomain(QString(domain));
}

int App_Exec() {
    return QGuiApplication::exec();
}
