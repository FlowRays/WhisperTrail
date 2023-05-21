package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/WhisperTrail/model"
	"github.com/FlowRays/WhisperTrail/service"
)

func GetImage(c *gin.Context, db *model.Database) {
	isLoggedIn_, exists := c.Get("isLoggedIn")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login state not found"})
		return
	}
	isLoggedIn, ok := isLoggedIn_.(bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "IsLoggedIn type error"})
		return
	}

	var userID uint
	if isLoggedIn {
		userID_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID not found"})
			return
		}
		userID, ok = userID_.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID type error"})
			return
		}
	}
	id_ := c.Param("id")

	id, err := strconv.ParseUint(id_, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ImageID error"})
		return
	}
	imageBytes, err := service.GetImage(uint(id), isLoggedIn, userID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
		return
	}

	// 设置响应头为图片类型
	c.Header("Content-Type", "image/jpeg")

	// 返回图片数据
	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}
