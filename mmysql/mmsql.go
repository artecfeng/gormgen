package mmysql

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

type Db struct {
	DB *gorm.DB
}

func (d *Db) GetDb() *gorm.DB {
	return d.DB
}

func (d *Db) Close() error {
	db, err := d.GetDb().DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func New(addr string) (Repo, error) {
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=%t&loc=%s",
		addr,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("连接数据库失败 %s", addr))
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	return &Db{
		DB: db,
	}, nil
}

type TableInfo struct {
	TableName    string         `db:"table_name"`    // name
	TableComment sql.NullString `db:"table_comment"` // comment
}

type TableColumn struct {
	OrdinalPosition uint16         `db:"ORDINAL_POSITION"` // position
	ColumnName      string         `db:"COLUMN_NAME"`      // name
	ColumnType      string         `db:"COLUMN_TYPE"`      // column_type
	DataType        string         `db:"DATA_TYPE"`        // data_type
	ColumnKey       sql.NullString `db:"COLUMN_KEY"`       // key
	IsNullable      string         `db:"IS_NULLABLE"`      // nullable
	Extra           sql.NullString `db:"EXTRA"`            // extra
	ColumnComment   sql.NullString `db:"COLUMN_COMMENT"`   // comment
	ColumnDefault   sql.NullString `db:"COLUMN_DEFAULT"`   // default value
}

func (d *Db) SelectTables(dbName string, tableName string) ([]TableInfo, error) {
	var tableSlice []TableInfo
	var tableNameSlice []string
	var commentSlice []sql.NullString

	sqlTables := fmt.Sprintf("SELECT `table_name`,`table_comment` FROM `information_schema`.`tables` WHERE `table_schema`= '%s'", dbName)
	rows, err := d.GetDb().Raw(sqlTables).Rows()
	if err != nil {
		return tableSlice, err
	}
	defer rows.Close()

	for rows.Next() {
		var info TableInfo
		err = rows.Scan(&info.TableName, &info.TableComment)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		tableSlice = append(tableSlice, info)
		tableNameSlice = append(tableNameSlice, info.TableName)
		commentSlice = append(commentSlice, info.TableComment)
	}

	if tableName != "*" {
		tableSlice = nil
		chooseTables := strings.Split(tableName, ",")
		indexMap := make(map[int]int)
		for _, item := range chooseTables {
			subIndexMap := getTargetIndexMap(tableNameSlice, item)
			for k, v := range subIndexMap {
				if _, ok := indexMap[k]; ok {
					continue
				}
				indexMap[k] = v
			}
		}

		if len(indexMap) != 0 {
			for _, v := range indexMap {
				var info TableInfo
				info.TableName = tableNameSlice[v]
				info.TableComment = commentSlice[v]
				tableSlice = append(tableSlice, info)
			}
		}
	}

	return tableSlice, err
}
func getTargetIndexMap(tableNameSlice []string, item string) map[int]int {
	indexMap := make(map[int]int, len(tableNameSlice))
	for i := 0; i < len(tableNameSlice); i++ {
		if tableNameSlice[i] == item {
			if _, ok := indexMap[i]; ok {
				continue
			}
			indexMap[i] = i
		}
	}
	return indexMap
}

func (d *Db) SelectTableColumn(dbName string, tableName string) ([]TableColumn, error) {
	var columns []TableColumn

	sqlTableColumn := fmt.Sprintf("SELECT `ORDINAL_POSITION`,`COLUMN_NAME`,`COLUMN_TYPE`,`DATA_TYPE`,`COLUMN_KEY`,`IS_NULLABLE`,`EXTRA`,`COLUMN_COMMENT`,`COLUMN_DEFAULT` FROM `information_schema`.`columns` WHERE `table_schema`= '%s' AND `table_name`= '%s' ORDER BY `ORDINAL_POSITION` ASC",
		dbName, tableName)

	rows, err := d.GetDb().Raw(sqlTableColumn).Rows()
	if err != nil {
		fmt.Printf("execute query table column action error, detail is [%v]\n", err.Error())
		return columns, err
	}
	defer rows.Close()

	for rows.Next() {
		var column TableColumn
		err = rows.Scan(
			&column.OrdinalPosition,
			&column.ColumnName,
			&column.ColumnType,
			&column.DataType,
			&column.ColumnKey,
			&column.IsNullable,
			&column.Extra,
			&column.ColumnComment,
			&column.ColumnDefault)
		if err != nil {
			fmt.Printf("query table column scan error, detail is [%v]\n", err.Error())
			return columns, err
		}
		columns = append(columns, column)
	}

	return columns, err
}
