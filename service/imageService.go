package service

import (
	"io/ioutil"

	"github.com/FlowRays/WhisperTrail/dao"
	"github.com/FlowRays/WhisperTrail/model"
)

func GetImage(id uint, isLoggedIn bool, userID uint, db *model.Database) ([]byte, error) {
	// 根据 ID 查询数据库，获取对应的路径
	var landmark model.Landmark
	landmark.ID = id
	err := dao.GetLandmarkByID(isLoggedIn, userID, &landmark, db)
	if err != nil {
		return nil, err
	}

	// 读取图片文件
	imageBytes, err := ioutil.ReadFile(landmark.Path)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}
