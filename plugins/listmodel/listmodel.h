#pragma once

#ifndef QAMEL_TABLEMODEL_H
#define QAMEL_TABLEMODEL_H

#ifdef __cplusplus

#include <QQmlEngine>
#include <QAbstractListModel>
#include <QVariant>
#include <QHash>
#include <QByteArray>
#include <QString>

class QamelListModel : public QAbstractListModel
{
    Q_OBJECT
    Q_PROPERTY(QVariantList contents READ contents WRITE setContents NOTIFY contentsChanged)

public:
    int rowCount(const QModelIndex & = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role) const override;
    QHash<int, QByteArray> roleNames() const override;
    bool insertRows(int, int, const QModelIndex & = QModelIndex()) override;
    bool removeRows(int, int, const QModelIndex & = QModelIndex()) override;
    bool moveRows(const QModelIndex &, int, int, const QModelIndex &, int) override;

    QVariantList contents() const;

public slots:
    void setContents(QVariantList);

    QVariant get(int);
    int count();
    void clear();
    void insertRow(int, QVariant);
    void insertRows(int, QVariantList);
    void appendRow(QVariant);
    void appendRows(QVariantList);
    void deleteRows(int, int);
    void setRow(int, QVariantMap);
    void setRowProperty(int, QString, QVariant);
    void swapRow(int, int);
    void moveRows(int, int, int);

signals:
    void contentsChanged();

private:
    QVariantList _contents;
};

extern "C" {
#endif // __cplusplus

// Register QML
void QamelListModel_RegisterQML(char* uri, int versionMajor, int versionMinor, char* qmlName);

#ifdef __cplusplus
}
#endif

#endif // QAMEL_TABLEMODEL_H
