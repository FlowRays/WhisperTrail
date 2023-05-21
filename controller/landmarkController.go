package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/WhisperTrail/model"
	"github.com/FlowRays/WhisperTrail/service"
)

func CreateLandmark(c *gin.Context, db *model.Database) {
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

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := "uploads/" + file.Filename
	latitude := c.PostForm("latitude")
	longitude := c.PostForm("longitude")
	text := c.PostForm("text")

	err = service.CreateLandmark(file, path, latitude, longitude, text, isLoggedIn, userID, c, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully!"})
}

func GetLandmark(c *gin.Context, db *model.Database) {
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

	landmarks, err := service.GetLandmark(isLoggedIn, userID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch landmarks"})
		return
	}
	c.JSON(http.StatusOK, landmarks)
}
