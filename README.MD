## 执行命令

```
git clone git@github.com:artecfeng/gormgen.git
cd gormgen
./gormgen.sh

//所有表
gormgen -addr "root:@tcp(localhost:3306)/ginadmin"
//单个或多个表
gormgen -addr "root:@tcp(localhost:3306)/ginadmin" -tables "admin"
//json支持
gormgen -addr "root:@tcp(localhost:3306)/ginadmin" -tables "admin" -json true
```

