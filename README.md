## go-gin 项目简单实战

#### 安装依赖

于 gogin 文件目录下，执行如下命令

```
go init go.mod

go mod tidy

go run main.go
```


#### 目录分布


- common 公共包
- conf：用于存储配置文件
- cron：定时器
- docs：文档
- middleware：应用中间件
- models：应用数据库模型
- node：笔记
- pkg：第三方包
- routers：路由逻辑处理
- runtime：应用运行时数据
- tem：html文件
- test：测试
#### 相关文档

- https://github.com/gin-gonic/gin

- https://github.com/kardianos/govendor （依赖包管理工具），参考链接： https://blog.csdn.net/bandaoyu/article/details/126414230

- https://github.com/go-ini/ini （配置文件）， https://ini.unknwon.io/ （中文文档）

- https://github.com/Unknwon/com （工具包）

- https://github.com/jinzhu/gorm

- https://github.com/go-sql-driver/mysql

- https://segmentfault.com/a/1190000013297625 （参考煎鱼🐟所写内容，以此文档为基础）

- https://github.com/astaxie/beego/tree/master/validation （beego的表单验证库）

- https://github.com/dgrijalva/jwt-go （身份验证）

- https://github.com/swaggo/swag （swag 文档）

- https://github.com/fvbock/endless （优雅关闭）

- https://github.com/robfig/cron （定时器）

- https://github.com/tealeg/xlsx ， https://github.com/qax-os/excelize

- https://github.com/boombuler/barcode （二维码）

- https://github.com/golang/freetype （图片绘制）

  