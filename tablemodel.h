#pragma once

#ifndef QAMEL_TABLEMODEL_H
#define QAMEL_TABLEMODEL_H

#ifdef __cplusplus

#include <QQmlEngine>
#include <QAbstractTableModel>
#include <QVariant>
#include <QHash>
#include <QByteArray>
#include <QString>

class QamelTableModel : public QAbstractTableModel
{
    Q_OBJECT
    Q_PROPERTY(int columns READ columns WRITE setColumns NOTIFY columnsChanged)
    Q_PROPERTY(QVariantList contents READ contents WRITE setContents NOTIFY contentsChanged)
    Q_PROPERTY(QVariantMap columnsName READ columnsName WRITE setColumnsName NOTIFY columnsNameChanged)

public:
    int rowCount(const QModelIndex & = QModelIndex()) const override;
    int columnCount(const QModelIndex & = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role) const override;
    QHash<int, QByteArray> roleNames() const override;
    bool insertRows(int, int, const QModelIndex & = QModelIndex()) override;
    bool removeRows(int, int, const QModelIndex & = QModelIndex()) override;
    bool moveRows(const QModelIndex &, int, int, const QModelIndex &, int) override;

    int columns() const;
    QVariantList contents() const;
    QVariantMap columnsName() const;

public slots:
    void setColumns(int);
    void setContents(QVariantList);
    void setColumnsName(QVariantMap columnsName);

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
    void columnsChanged();
    void contentsChanged();
    void columnsNameChanged();

private:
    int _columns;
    QVariantList _contents;
    QVariantMap _columnsName;

    QString getColumnName(int) const;
};

extern "C" {
#endif // __cplusplus

// Register QML
void QamelTableModel_RegisterQML(char* uri, int versionMajor, int versionMinor, char* qmlName);

#ifdef __cplusplus
}
#endif

#endif // QAMEL_TABLEMODEL_H
