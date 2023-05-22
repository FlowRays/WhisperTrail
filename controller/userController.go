package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/WhisperTrail/model"
	"github.com/FlowRays/WhisperTrail/service"
)

func UserRegister(c *gin.Context, db *model.Database) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.UserRegister(&user, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func UserLogin(c *gin.Context, db *model.Database) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := service.UserLogin(&user, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.Username, "token": token})
}
