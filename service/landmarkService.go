package service

import (
	"mime/multipart"

	"github.com/FlowRays/WhisperTrail/dao"
	"github.com/FlowRays/WhisperTrail/model"
	"github.com/gin-gonic/gin"
)

func GetLandmark(isLoggedIn bool, userID uint, db *model.Database) ([]model.Landmark, error) {
	var landmarks []model.Landmark
	err := dao.GetLandmarksByUserStatus(isLoggedIn, userID, &landmarks, db)

	return landmarks, err
}

func CreateLandmark(file *multipart.FileHeader, path string, latitude string, longitude string, text string, isLoggedIn bool, userID uint, c *gin.Context, db *model.Database) error {
	// 将上传的文件保存到本地
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		return err
	}
	landmark := model.Landmark{Path: path, Latitude: latitude, Longitude: longitude, Text: text, IsLoggedIn: isLoggedIn, UserID: userID}
	err = dao.CreateLandmark(&landmark, db)

	return err
}
