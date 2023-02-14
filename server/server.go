package server

import (
	"fmt"
	"gormgen/mmysql"
	"log"
	"os"
	"strings"
)

func DoGen(addr, dbName, tableName string, withjson bool) error {
	db, err := mmysql.New(addr)
	if err != nil {
		log.Fatal("new db err", err)
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db close err", err)
		}
	}()

	tables, err := db.SelectTables(dbName, tableName)
	if err != nil {
		return err
	}
	for _, table := range tables {

		filepath := "./model/" + table.TableName
		os.MkdirAll(filepath, 0766)
		fmt.Println("create dir : ", filepath)

		modelName := fmt.Sprintf("%s/%s_model.go", filepath, table.TableName)
		modelFile, err := os.OpenFile(modelName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
		if err != nil {
			fmt.Printf("create and open model file error %v\n", err.Error())
			return err
		}
		fmt.Println("  └── file : ", table.TableName+"_model.go")
		columnInfo, err := db.SelectTableColumn(dbName, table.TableName)

		modelContent := fmt.Sprintf("package %s\n", table.TableName)
		for _, v := range columnInfo {
			if textType(v.DataType) == "time.Time" {
				modelContent += fmt.Sprintf(`import "time"`)
				break
			}
		}
		modelContent += fmt.Sprintf("\n\n// %s %s \n", Case2Camel(table.TableName), table.TableComment.String)
		modelContent += fmt.Sprintf("//go:generate gormgen -structs %s -input . \n", Case2Camel(table.TableName))
		modelContent += fmt.Sprintf("type %s struct {\n", Case2Camel(table.TableName))

		if err != nil {
			continue
		}
		for _, info := range columnInfo {
			if withjson {
				modelContent += fmt.Sprintf("	%s %s `json:\"%s\" gorm:\"%s\"` // %s\n", Case2Camel(info.ColumnName), textType(info.DataType), info.ColumnName, info.ColumnName, info.ColumnComment.String)
			} else {
				modelContent += fmt.Sprintf("	%s %s `gorm:\"%s\"` // %s\n", Case2Camel(info.ColumnName), textType(info.DataType), info.ColumnName, info.ColumnComment.String)
			}
		}
		modelContent += "}\n"
		modelFile.WriteString(modelContent)

		modelFile.Close()

	}

	return nil
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}
func textType(s string) string {
	var mysqlTypeToGoType = map[string]string{
		"tinyint":    "int32",
		"smallint":   "int32",
		"mediumint":  "int32",
		"int":        "int32",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float64",
		"double":     "float64",
		"decimal":    "float64",
		"date":       "string",
		"time":       "string",
		"year":       "string",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"char":       "string",
		"varchar":    "string",
		"tinyblob":   "string",
		"tinytext":   "string",
		"blob":       "string",
		"text":       "string",
		"mediumblob": "string",
		"mediumtext": "string",
		"longblob":   "string",
		"longtext":   "string",
	}
	return mysqlTypeToGoType[s]
}
