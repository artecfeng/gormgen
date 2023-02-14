## 执行命令

```
go build main.go
//添加到bin目录，可全局使用
cp main $GOPATH/bin/gormgen

//所有表
gormgen -addr "root:@tcp(localhost:3306)/ginadmin"
//单个或多个表
gormgen -addr "root:@tcp(localhost:3306)/ginadmin" -tables "admin"
//json支持
gormgen -addr "root:@tcp(localhost:3306)/ginadmin" -tables "admin" -json true
```

