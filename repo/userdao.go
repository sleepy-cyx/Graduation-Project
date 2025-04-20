package repo

import (
	"Graduation-Project/model"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	DB *gorm.DB
}

// GetUserInfoByUsername 根据username查表
func (usr *UserDao) GetUserInfoByUsername(userId uint) (*model.User, bool, error) {
	user := model.User{}
	result := usr.DB.Where("user_name = ?", username).Limit(1).Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 没有找到记录
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &user, true, nil

}

// AlterUserInfo 修改用户信息
func (usr *UserDao) AlterUserInfo(username string, nickname string, pictureUrl string, loginIP string, loginTime time.Time) (bool, error) {
	tableName := getTableName(username, 10)
	updateTime := time.Now().Format(time.RFC3339)
	updateTimeParsed, err := time.Parse(time.RFC3339, updateTime)
	if err != nil {
		return false, err
	}
	// 创建一个map来存储要更新的字段
	updates := make(map[string]interface{})
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if pictureUrl != "" {
		updates["picture_url"] = pictureUrl
	}
	if loginIP != "" {
		updates["latest_ip"] = loginIP
	}
	// LoginTime总是需要更新
	updates["latest_login"] = loginTime
	updates["update_time"] = updateTimeParsed
	if nickname == "" {
		return true, nil
	}
	// 使用Where找到特定的用户并更新字段
	result := usr.DB.Table(tableName).Where("user_name = ?", username).Limit(1).Updates(updates)
	// 检查是否有错误发生
	if result.Error != nil {
		return false, result.Error
	}
	// 检查是否找到了记录
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
func getTableName(userName string, numTables int) string {
	hash := sha256.Sum256([]byte(userName))
	hashNum := binary.BigEndian.Uint64(hash[:8])
	tableIdx := hashNum % uint64(numTables)
	tableName := fmt.Sprintf("user_tab_0000000%d", tableIdx)
	return tableName
}
