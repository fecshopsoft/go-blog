package model

import (    
    "fmt"
    "strconv"
    "github.com/gin-gonic/gin"  
    mysqlPool "github.com/fecshopsoft/golang-db/mysql"
)

type customer struct {
    Id      int    `form:"id" json:"id" `
    Name    string `form:"name" json:"name" binding:"required"`
    Age     int    `form:"age" json:"age" binding:"required"`
}

type CustomerData struct{}

var Customer CustomerData

func (Customer CustomerData) List(mysqlDB *mysqlPool.SQLConnPool) gin.H{
    body := make(gin.H) 
    rows, err := mysqlDB.Query("SELECT * From customer")
    if err != nil {
        fmt.Printf("%s\r\n","mysql query error")
    }
    //fmt.Printf("%v\r\n",rows)
    var dbdata []gin.H
    if rows != nil {
        for _, row := range rows {
            dbdata = append(dbdata, gin.H(row))
        }
    }
    body["status"] = 200
    body["data"] = dbdata
    return body
}

func (Customer CustomerData) AddOne(mysqlDB *mysqlPool.SQLConnPool, c *gin.Context) gin.H{
    body := make(gin.H) 
    // 保存
    var json customer
    if err := c.ShouldBindJSON(&json); err == nil {
        lastId, err := mysqlDB.Insert("INSERT INTO customer (`name`, `age`) VALUES( ?, ? )", json.Name, json.Age) // ? = placeholder
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
        body["insertId"] = lastId
        body["status"] = "success"
    } else {
        body["status"] = err.Error()
    }
    return  body
}

func (Customer CustomerData) UpdateById(mysqlDB *mysqlPool.SQLConnPool, c *gin.Context) gin.H{
    userId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        panic("userId can not empty")  
    }
    body := make(gin.H) 
    // 保存
    var json customer
    if err := c.ShouldBindJSON(&json); err == nil {
        // 进行数据库操作
        affect, err := mysqlDB.Update("update customer set `name` = ? , `age` = ? where `id` = ? ", json.Name, json.Age, userId) // ? = placeholder
        if err != nil {
            panic(err.Error()) 
        }
        body["updateCount"] = affect
        body["status"] = "success"
    } else {
        body["status"] = err.Error()
    }
    return  body
}

func (Customer CustomerData) DeleteById(mysqlDB *mysqlPool.SQLConnPool, c *gin.Context) gin.H{
    userId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        panic("userId can not empty")  
    }
    body := make(gin.H) 
    affect, err := mysqlDB.Update("delete from customer where `id` = ?", userId) // ? = placeholder
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    body["deleteCount"] = affect
    return  body
}


func (Customer CustomerData) Transaction(mysqlDB *mysqlPool.SQLConnPool, c *gin.Context) gin.H{
    userId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        panic("userId can not empty")  
    }
    body := make(gin.H) 
    
    sqlransaction, err := mysqlDB.Begin()
    if err != nil {
        panic("transaction begin error")  
    }
    
    affect1, err := sqlransaction.Update("update customer set `name` = ? , `age` = ? where `id` = ? ", "111",2, userId) // ? = placeholder
    if err != nil {
        sqlransaction.Rollback() 
    }
    
    affect2, err := sqlransaction.Update("update customer set `name` = ? , `age` = ? where `id` = ? ", "222", 2, userId) // ? = placeholder
    if err != nil {
        sqlransaction.Rollback() 
    }
    
    lastId, err := sqlransaction.Insert("INSERT INTO customer (`name`, `age`) VALUES( ?, ? )", "333333", 333) // ? = placeholder
    if err != nil {
        sqlransaction.Rollback()
    }
        
    sqlransaction.Commit()
    body["affect1"] = affect1
    body["affect2"] = affect2
    body["lastId"]  = lastId
    return  body
}