package main

import (
	"flag"
	"fmt"
	"gormgen/mfmt"
	"gormgen/server"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", "root:@tcp(127.0.0.1:3306)/ginadmin", "请输入 db 地址，例如：root:@tcp(127.0.0.1:3306)/ginadmin\n")
	table := flag.String("tables", "*", "请输入 table 名称，默认为“*”，多个可用“,”分割\n")
	flag.Parse()
	fmt.Println(*addr)
	fmt.Println(*table)
	split := strings.Split(*addr, "/")
	if len(split) != 2 {
		os.Exit(0)
	}
	dbName := split[1]
	dbTables := strings.ToLower(*table)
	err := server.DoGen(*addr, dbName, dbTables)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("create model success")
	fmt.Println("start generate...")
	mfmt.DoFmt()
}
