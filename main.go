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

type Landmark struct {
	ID        uint   `gorm:"primaryKey"`
	Path      string `gorm:"size:255"`
	Latitude  string `gorm:"not null"`
	Longitude string `gorm:"not null"`
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
	err = db.AutoMigrate(&Landmark{})
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
	// 定义处理 CORS 请求的中间件
	// corsMiddleware := func() gin.HandlerFunc {
	// 	return func(c *gin.Context) {
	// 		// 允许的域名列表
	// 		allowedOrigins := []string{"http://localhost:3000", "http://172.27.144.70:5941"}

	// 		origin := c.Request.Header.Get("Origin")
	// 		// 检查请求的来源是否在允许的域名列表中
	// 		for _, allowedOrigin := range allowedOrigins {
	// 			if allowedOrigin == origin {
	// 				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	// 				break
	// 			}
	// 		}

	// 		// 设置其他 CORS 相关的响应头
	// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 		// 处理预检请求
	// 		if c.Request.Method == "OPTIONS" {
	// 			c.AbortWithStatus(200)
	// 			return
	// 		}

	// 		// 调用下一个中间件或处理函数
	// 		c.Next()
	// 	}
	// }

	// 应用 CORS 中间件
	// r.Use(corsMiddleware())

	r.POST("/upload", func(c *gin.Context) {
		createLandmark(c, db)
	})

	r.GET("/get", func(c *gin.Context) {
		getAllLandmark(c, db)
	})

	r.GET("/image/:id", func(c *gin.Context) {
		getImage(c, db)
	})
}

func createLandmark(c *gin.Context, db *Database) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	latitude := c.PostForm("latitude")
	longitude := c.PostForm("longitude")

	path := "uploads/" + file.Filename

	// 将上传的文件保存到本地
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	landmark := Landmark{Path: path, Latitude: latitude, Longitude: longitude}
	result := db.DB.Create(&landmark)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully!"})
}

func getAllLandmark(c *gin.Context, db *Database) {
	var landmarks []Landmark
	if err := db.DB.Find(&landmarks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch landmarks"})
		return
	}
	c.JSON(http.StatusOK, landmarks)
}

func getImage(c *gin.Context, db *Database) {
	id := c.Param("id")

	// 根据 ID 查询数据库，获取对应的路径
	var landmark Landmark
	if err := db.DB.Where("id = ?", id).First(&landmark).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// 读取图片文件
	imageBytes, err := ioutil.ReadFile(landmark.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
		return
	}

	// 设置响应头为图片类型
	c.Header("Content-Type", "image/jpeg")

	// 返回图片数据
	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}
