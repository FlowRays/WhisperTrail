package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/FlowRays/WhisperTrail/controller"
	"github.com/FlowRays/WhisperTrail/model"
)

func main() {
	err := readConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据库迁移
	err = db.AutoMigrate(&model.Landmark{}, &model.User{}, &model.Rate{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建数据库实例
	database := model.Database{
		DB: db,
	}

	r := gin.Default()
	registerAPI(r, &database)
	r.Run(":8080")
}

func readConfig() error {
	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	return err
}

func connectDB() (*gorm.DB, error) {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func registerAPI(r *gin.Engine, db *model.Database) {
	r.POST("/api/register", func(c *gin.Context) {
		controller.UserRegister(c, db)
	})

	r.POST("/api/login", func(c *gin.Context) {
		controller.UserLogin(c, db)
	})

	auth := r.Group("/api")
	auth.Use(controller.AuthMiddleware())

	auth.POST("/upload", func(c *gin.Context) {
		controller.CreateLandmark(c, db)
	})

	auth.GET("/get", func(c *gin.Context) {
		controller.GetLandmark(c, db)
	})

	auth.GET("/image/:id", func(c *gin.Context) {
		controller.GetImage(c, db)
	})

	r.POST("/api/rate", func(c *gin.Context) {
		controller.CreateRate(c, db)
	})

	r.GET("/api/rate/:id", func(c *gin.Context) {
		controller.GetRate(c, db)
	})
}
