package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/WhisperTrail/model"
	"github.com/FlowRays/WhisperTrail/service"
)

func CreateRate(c *gin.Context, db *model.Database) {
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
	if !isLoggedIn {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	var rate model.Rate
	if err := c.ShouldBind(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateRate(&rate, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Create rate successfully"})
}

func GetRate(c *gin.Context, db *model.Database) {
	id_ := c.Param("id")
	id, err := strconv.ParseUint(id_, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LandmarkID error"})
		return
	}

	rating, err := service.GetRate(uint(id), db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rating"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get rating successfully", "rating": rating})
}
