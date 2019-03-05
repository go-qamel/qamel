#include "quickstyle.h"
#include <QQuickStyle>
#include <QString>

void SetQuickStyle(char* style) {
    QQuickStyle::setStyle(QString(style));
}

void SetQuickFallbackStyle(char* style) {
    QQuickStyle::setFallbackStyle(QString(style));
}

void AddQuickStylePath(char* style) {
    QQuickStyle::addStylePath(QString(style));
}
