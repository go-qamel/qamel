#include "_cgo_export.h"
#include "tablemodel.h"

// Overrided method for QAbstractTableModel
int QamelTableModel::columnCount(const QModelIndex &) const {
    return _columns;
}

int QamelTableModel::rowCount(const QModelIndex &) const {
    return _contents.size();
}

QVariant QamelTableModel::data(const QModelIndex &index, int role) const {
    if (role != Qt::DisplayRole) {
        return QVariant();
    }

    QString colName = getColumnName(index.column());
    QVariant obj = _contents.at(index.row());
    return obj.toMap().value(colName);
}

QHash<int, QByteArray> QamelTableModel::roleNames() const {
    QHash<int, QByteArray> role;
    role[Qt::DisplayRole] = "display";
    return role;
}

bool QamelTableModel::insertRows(int row, int count, const QModelIndex &parent) {
    return QAbstractTableModel::insertRows(row, count, parent);
}

bool QamelTableModel::removeRows(int row, int count, const QModelIndex &parent) {
    return QAbstractTableModel::removeRows(row, count, parent);
}

bool QamelTableModel::moveRows(const QModelIndex &sourceParent, int sourceRow, int count, const QModelIndex &destinationParent, int destinationChild) {
    return QAbstractTableModel::moveRows(sourceParent, sourceRow, count, destinationParent, destinationChild);
}

// Methods for custom QML properties
int QamelTableModel::columns() const {
    return _columns;
}

QVariantList QamelTableModel::contents() const {
    return _contents;
}

QVariantMap QamelTableModel::columnsName() const {
    return _columnsName;
}

void QamelTableModel::setColumns(int columns) {
    _columns = columns;
    emit columnsChanged();
    emit layoutChanged();
}

void QamelTableModel::setContents(QVariantList contents) {
    beginResetModel();
    _contents = contents;
    endResetModel();

    emit contentsChanged();
}

void QamelTableModel::setColumnsName(QVariantMap columnsName) {
    _columnsName = columnsName;
    emit columnsNameChanged();
    emit layoutChanged();
}

// Private method
QString QamelTableModel::getColumnName(int column) const {
    QVariant val = _columnsName[QString::number(column)];
    return val.toString();
}

// Methods for custom slots
QVariant QamelTableModel::get(int row) {
    return _contents.at(row);
}

int QamelTableModel::count() {
    return _contents.count();
}

void QamelTableModel::clear() {
    beginResetModel();
    _contents.clear();
    endResetModel();
}

void QamelTableModel::insertRow(int row, QVariant obj) {
    beginInsertRows(QModelIndex(), row, row);
    _contents.insert(row, obj);
    endInsertRows();

    emit contentsChanged();
}

void QamelTableModel::insertRows(int row, QVariantList objects) {
    beginInsertRows(QModelIndex(), row, row+objects.count()-1);
    for (int i = 0; i < objects.count(); ++i) {
        _contents.insert(row+i, objects[i]);
    }
    endInsertRows();

    emit contentsChanged();
}

void QamelTableModel::appendRow(QVariant obj) {
    int row = _contents.count();

    beginInsertRows(QModelIndex(), row, row);
    _contents.insert(row, obj);
    endInsertRows();

    emit contentsChanged();
}

void QamelTableModel::appendRows(QVariantList objects) {
    int row = _contents.count();
    int last = row + objects.count() - 1;

    beginInsertRows(QModelIndex(), row, last);
    for (int i = 0; i < objects.count(); ++i) {
        _contents.insert(row+i, objects[i]);
    }
    endInsertRows();

    emit contentsChanged();
}

void QamelTableModel::deleteRows(int row, int count) {
    if (row < 0 || row >= rowCount()) return;

    int last = row+count-1;

    beginRemoveRows(QModelIndex(), row, last);
    for (int i = last; i >= row; --i) {
        _contents.removeAt(i);
    }
    endRemoveRows();

    emit contentsChanged();
}

void QamelTableModel::setRow(int row, QVariantMap newObj) {
    if (row < 0 || row >= rowCount()) return;

    _contents[row] = newObj;
    emit contentsChanged();
    emit dataChanged(createIndex(row, 0), createIndex(row, columnCount()-1));
}

void QamelTableModel::setRowProperty(int row, QString propName, QVariant newProp) {
    if (row < 0 || row >= rowCount()) return;

    QVariant obj = _contents.at(row);
    QVariantMap map = obj.toMap();
    map[propName] = newProp;
    _contents[row] = QVariant(map);

    emit contentsChanged();
    emit dataChanged(createIndex(row, 0), createIndex(row, columnCount()-1));
}

void QamelTableModel::swapRow(int i, int j) {
    if (i == j) return;
    if (i < 0 || i >= rowCount()) return;
    if (j < 0 || j >= rowCount()) return;

    _contents.swap(i, j);
    emit dataChanged(createIndex(i, 0), createIndex(i, columnCount()-1));
    emit dataChanged(createIndex(j, 0), createIndex(j, columnCount()-1));
}

void QamelTableModel::moveRows(int first, int last, int dst) {
    int lastDst = dst + (last - first);

    if (first < 0 || first >= rowCount()) return;
    if (last < 0 || last >= rowCount()) return;
    if (dst< 0 || dst >= rowCount()) return;
    if (lastDst < 0 || lastDst >= rowCount()) return;

    beginMoveRows(QModelIndex(), first, last, QModelIndex(), dst);
    for (int i = last; i >= first; --i) {
        _contents.move(i, lastDst - (last - i));
    }
    endMoveRows();

    emit contentsChanged();
}

void QamelTableModel_RegisterQML(char* uri, int versionMajor, int versionMinor, char* qmlName) {
	qmlRegisterType<QamelTableModel>(uri, versionMajor, versionMinor, qmlName);
}

#include "moc-tablemodel.h"