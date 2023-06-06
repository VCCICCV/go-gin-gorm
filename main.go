package main

// PROJECT_NAME:go-gin-gorm
// DATE:2023/6/5 21:24
// USER:Administrator
// 错误码 正确200 无效|错误400
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
		Phone int `gorm:"type:int;not null"json:"phone" binding:"required"`
		Email string `gorm:"type:varchar(40);not null"json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200);not null"json:"address" binding:"required"`
	}

	db.AutoMigrate(&List{})

	// 接口
	r := gin.Default()

	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "请求成功",
	//	})
	//})

	// create
	r.POST("/user/add",func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)
		fmt.Println(err)
		if err != nil{
			c.JSON(200, gin.H{
				"msg":"添加失败",
				"data":gin.H{},
				"code":400,
			})
		}else{
			// 操作数据库
			db.Create(&data) // 创建一条数据
			c.JSON(200, gin.H{
				"msg":"添加成功",
				"data":data,
				"code":200,
			})
		}
	})
	// delete
	r.DELETE("/user/delete/:id",func(c *gin.Context) {
		var data []List
		// 接收id
		id := c.Param("id")
		// 判断是否存在
		db.Where("id = ?",id).Find(&data)
		// id存在的情况,则删除,不存在则报告
		if len(data) == 0{
			c.JSON(200,gin.H{
				"mag":"id没有找到",
				"code":400,
			})
		}else{
			// 操作删除
			db.Where("id = ?",id).Delete(&data)

            c.JSON(200,gin.H{
                "msg":"删除成功",
                "code":200,
            })
		}
	})
	// update
	r.PUT("/user/update/:id",func(c *gin.Context) {
		var data List
		// 接收id
		id := c.Param("id")
		// 判断id是否存在
		db.Select("id").Where("id = ?",id).Find(&data)

		if data.ID == 0{
			c.JSON(200,gin.H{
                "mag":"id没有找到",
                "code":400,
            })
		}else{
			err := c.ShouldBindJSON(&data)

			if err!= nil{
				c.JSON(200,gin.H{
                    "msg":"更新失败",
                    "data":gin.H{},
                    "code":400,
                })
			}else{
				// 修改数据库内容
				db.Where("id = ?", id).Updates(&data)

				c.JSON(200,gin.H{
					"msg":"更新成功",
					"data":data,
					"code":200,
				})
			}

		}
	})
	// read 条件查询|全部查询|分页查询
	r.GET("/user/list/:name" ,func(c *gin.Context){
		name := c.Param("name")
		var dataList []List
		//	查询数据库
		db.Where("name = ?",name).Find(&dataList)
	    // 判断是否查询到数据
		if len(dataList) == 0{
			c.JSON(200,gin.H{
                "msg":"查询失败",
                "data":gin.H{},
                "code":400,
            })
		}else{
			c.JSON(200,gin.H{
                "msg":"查询成功",
                "data":dataList,
                "code":200,
            })
		}
	})
	// 端口号
	r.Run(":3004")
}
