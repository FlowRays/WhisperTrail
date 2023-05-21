package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type Database struct {
	DB *gorm.DB
}

type Image struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `gorm:"size:255"`
}

func main() {
	// 读取配置文件
	config, err := readConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 连接数据库
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据库迁移
	err = db.AutoMigrate(&Image{})
	if err != nil {
		log.Fatal(err)
	}

	// 创建数据库实例
	database := Database{
		DB: db,
	}

	r := gin.Default()
	registerAPI(r, &database)
	r.Run(":8080")
}

func readConfig(filename string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("Failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("Failed to unmarshal config data: %v", err)
	}

	return config, nil
}

func connectDB(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func registerAPI(r *gin.Engine, db *Database) {
	r.POST("/upload", func(c *gin.Context) {
		uploadImage(c, db)
	})
}

func uploadImage(c *gin.Context, db *Database) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := "uploads/" + file.Filename

	// 将上传的文件保存到本地
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	image := Image{Path: path}
	result := db.DB.Create(&image)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully!"})
}
