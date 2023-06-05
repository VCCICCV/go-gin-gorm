package main

// PROJECT_NAME:go-gin-gorm
// DATE:2023/6/5 21:24
// USER:Administrator
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)
func main(){
	//连接数库
	dsn := "root:666666@tcp(127.0.0.1:3306)/go-crud-demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 建表不复数
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()
	// SetMaxIdleConns 空闲连接池最大连接数
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns s打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 连接可服用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 结构体
	type List struct {
		// 添加主键
		gorm.Model
		Name string `gorm:"type:varchar(20);not null"json:"name" binding:"required"`
		Status string `gorm:"type:varchar(20);not null"json:"status" binding:"required"`
		Phone string `gorm:"type:int;not null"json:"phone" binding:"required"`
		Email string `gorm:"type:varchar(40);not null"json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200);not null"json:"address" binding:"required"`
	}

	db.AutoMigrate(&List{})

	// 接口
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// create

	// delete

	// update

	// read



	// 端口号
	PORT := "3004"
	r.Run(":" + PORT)
}