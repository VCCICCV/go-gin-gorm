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
	"strconv"
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
	r.GET("/user/list" ,func(c *gin.Context){
		var dataList []List
        // 查询全部数据，分页数据
		pageNum,_ := strconv.Atoi(c.Query("pageNum"))
		pageSize,_ := strconv.Atoi(c.Query("pageSize"))
		//println(pageNum)
		//println(pageSize)

		// 判断是否需要分页
		if pageNum == 0{
            pageNum = -1
        }
		if pageSize == 0{
			pageSize = -1
		}
		// 需要跳过的记录数
		offsetVal := (pageNum - 1) * pageSize
		// 不进行分页查询
		if pageNum == -1 && pageSize == -1{
			offsetVal = -1
		}
		// 查询数据库
		var total int64
		//db.Model(dataList)：指定查询操作的数据模型为 dataList，即查询 dataList 对应的数据表。
		//Count(&total)：查询数据表中符合条件的记录数，并将结果保存到变量 total 中。该函数的参数是一个指针类型，用于接收查询结果。
		//Limit(-1)：设置查询结果的最大记录数为 -1，表示不限制查询结果的记录数。
		//Offset(-1)：设置查询结果的偏移量为 -1，表示从倒数第一个记录开始查询。
		//Find(&dataList)：执行查询操作，并将查询结果保存到变量 dataList 中。该函数的参数是一个指针类型，用于接收查询结果。
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList)==0{
			c.JSON(200,gin.H{
                "msg":"查询失败",
                "data":gin.H{},
                "code":400,
            })
		}else{
			c.JSON(200,gin.H{
                "msg":"查询成功",
                "data":gin.H{
					"list":dataList,
					"total":total,
					"pageNum":pageNum,
                    "pageSize":pageSize,
				},
                "code":200,
            })
		}
	})
	// 端口号
	r.Run(":3004")
}
