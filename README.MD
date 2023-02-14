## 执行命令

在根目录下执行脚本：`./gormgen.sh `；

- addr：数据库地址，例如："root:@tcp(localhost:3306)/ginadmin"
- tables：表名，默认为 *，多个表名可用“,”分割，例如："admin"

USE：

```
go build main.go
cp main $GOPATH/bin/gormgen
gormgen -addr "root:@tcp(localhost:3306)/ginadmin" -tables "admin"
```

