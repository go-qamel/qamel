#include "application.h"
#include <QApplication>
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

    QApplication::setAttribute(Qt::AA_EnableHighDpiScaling);
    return new QApplication(argcs, argvs);
}

void App_SetAttribute(long long attribute, bool on) {
    QApplication::setAttribute(Qt::ApplicationAttribute(attribute), on);
}

void App_SetFont(char *family, int pointSize, int weight, bool italic) {
    QFont font = QFont(QString(family), pointSize, weight, italic);
    QApplication::setFont(font);
}

void App_SetQuitOnLastWindowClosed(bool quit) {
    QApplication::setQuitOnLastWindowClosed(quit);
}

void App_SetApplicationDisplayName(char* name) {
	QApplication::setApplicationDisplayName(QString(name));
}

void App_SetWindowIcon(char* fileName) {
    QIcon icon = QIcon(QString(fileName));
    QApplication::setWindowIcon(icon);
}

void App_SetApplicationName(char* name) {
    QApplication::setApplicationName(QString(name));
}

void App_SetApplicationVersion(char* version) {
    QApplication::setApplicationVersion(QString(version));
}

void App_SetOrganizationName(char* name) {
    QApplication::setOrganizationName(QString(name));
}

void App_SetOrganizationDomain(char* domain) {
    QApplication::setOrganizationDomain(QString(domain));
}

int App_Exec() {
    return QApplication::exec();
}