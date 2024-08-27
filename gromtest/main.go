package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func OpenDB() (*gorm.DB, error) {
	dsn := "root:yjfc4883212@tcp(127.0.0.1:3306)/xintest?charset=utf8mb4&parseTime=True&loc=Local"
	// 设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出的位置
		logger.Config{
			SlowThreshold: time.Second, // 慢sql的阀值
			LogLevel:      logger.Info, // Log level ；Silent、Error、Warn、Info；info 表示所有sql都会打印
			//IgnoreRecordNotFoundError: true,          // 忽略记录器的 ErrRecordNotFound 错误
			//ParameterizedQueries:      true,          // 不要在 SQL 日志中包含参数
			Colorful: true, // 是否禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, //设置全局的日志级别
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: "xin_",
		//},
	})
	if err != nil {
		panic(err)
	}

	return db, err
}

type Language struct {
	gorm.Model
	Name string
}

// 修改单个表名
func (Language) TableName() string {
	return "xin_test_language"
}

func main() {
	db, err := OpenDB()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Language{})

}
