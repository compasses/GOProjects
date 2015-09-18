# 说明
通过一个店名或者店的名字获取该店对应的schema name。

##Usage:
```
  -h string
        DB host name (default "www.xxx.com:3306")
  -p string
        DB password (default "12345")
  -q int
        eshop tenant id (default -1)
  -sn string
        eshop name
  -u string
        DB user name (default "root")
```

使用了[go-sql-driver](http://godoc.org/github.com/go-sql-driver/mysql)
简单易用可以用来开发更多的使用工具。数据库操作的大多数需求都能完成。