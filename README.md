# go-blog
go-blog



```
package main

import(
    "github.com/gin-gonic/gin"
    "net/http" 
    _ "github.com/go-sql-driver/mysql" 
    mysqlPool "github.com/fecshopsoft/golang-db/mysql"
    "github.com/fecshopsoft/go-blog/model"
)

func mysqlDBPool() *mysqlPool.SQLConnPool{
    host := `127.0.0.1:3306`
    database := `go_test`
    user := `root`
    password := `xxxx`
    charset := `utf8`
    // 用于设置最大打开的连接数
    maxOpenConns := 200
    // 用于设置闲置的连接数
    maxIdleConns := 100
    mysqlDB := mysqlPool.InitMySQLPool(host, database, user, password, charset, maxOpenConns, maxIdleConns)
    return mysqlDB
}

func main() { 
    mysqlDB := mysqlDBPool();
	r := gin.Default()
    v2 := r.Group("/v2")
    {
        // 查询部分
        v2.GET("/customers", func(c *gin.Context) {
            data := model.Customer.List(mysqlDB);
            c.JSON(http.StatusOK, data)
        })
        v2.POST("/customers", func(c *gin.Context) {
            data := model.Customer.AddOne(mysqlDB, c);
            c.JSON(http.StatusOK, data)
        })
        v2.PATCH("/customers/:id", func(c *gin.Context) {
            data := model.Customer.UpdateById(mysqlDB, c);
            c.JSON(http.StatusOK, data)
        })
        v2.DELETE("/customers/:id", func(c *gin.Context) {
            data := model.Customer.DeleteById(mysqlDB, c);
            c.JSON(http.StatusOK, data)
        })
        
        v2.POST("/customers/transaction/:id", func(c *gin.Context) {
            data := model.Customer.Transaction(mysqlDB, c);
            c.JSON(http.StatusOK, data)
        })
    }
    r.Run("120.24.37.249:3000") // 这里改成您的ip和端口
}
```