package mmysql

import "gorm.io/gorm"

type Repo interface {
	GetDb() *gorm.DB
	Close() error
	SelectTables(dbName string, tableName string) ([]TableInfo, error)
	SelectTableColumn(dbName string, tableName string) ([]TableColumn, error)
}
