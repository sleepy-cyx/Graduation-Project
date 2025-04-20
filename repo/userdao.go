package repo

import (
	"Graduation-Project/model"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDao struct {
	DB *gorm.DB
}

// 修改返回类型为切片（非指针）
func (usr *UserDao) GetScheduleInfoByUserId(userId uint) ([]model.Schedule, bool, error) {
	var schedules []model.Schedule
	result := usr.DB.Where("user_id = ?", userId).Find(&schedules)

	if result.Error != nil {
		return nil, false, result.Error
	}

	if len(schedules) == 0 {
		return schedules, false, nil // 直接返回空切片
	}

	return schedules, true, nil
}

// GetUserInfoByUsername 根据username查表
func (usr *UserDao) GetUserInfoByUsername(username string) (*model.User, bool, error) {
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

// CreateUser 创建新用户（注册逻辑）
func (usr *UserDao) CreateUser(user *model.User) (*model.User, bool, error) {
	// 执行数据库插入操作
	result := usr.DB.Create(user)

	if result.Error != nil {
		// 处理唯一约束冲突（用户名重复）
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, false, fmt.Errorf("用户名 %s 已存在", user.Username)
		}

		// 其他数据库错误
		return nil, false, result.Error
	}

	// 插入成功时返回用户信息（包含自动生成的ID）
	return user, true, nil
}

// CreateSchedule 创建新日程
func (dao *UserDao) CreateSchedule(schedule *model.Schedule) (*model.Schedule, bool, error) {
	// 检查外键约束（确保UserID存在）
	var user model.User
	if err := dao.DB.First(&user, schedule.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, fmt.Errorf("用户ID %d 不存在", schedule.UserID)
		}
		return nil, false, err
	}

	// 执行插入操作
	result := dao.DB.Create(schedule)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062: // 唯一约束冲突（假设有其他唯一字段）
				return nil, false, fmt.Errorf("日程冲突: %v", mysqlErr.Message)
			case 1452: // 外键约束失败（理论上已提前检查，此处为兜底）
				return nil, false, fmt.Errorf("关联用户不存在")
			}
		}
		return nil, false, result.Error
	}

	return schedule, true, nil
}

// UpdateSchedule 更新日程（根据ID和UserID验证权限）
func (dao *UserDao) UpdateSchedule(schedule *model.Schedule) (*model.Schedule, bool, error) {
	// 查询现有数据（防止越权修改）
	existing := model.Schedule{}
	result := dao.DB.Where("id = ? AND user_id = ?", schedule.Id, schedule.UserID).First(&existing)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, fmt.Errorf("日程不存在或无权修改")
		}
		return nil, false, result.Error
	}

	// 执行更新（避免更新UserID）
	updateData := map[string]interface{}{
		"title":      schedule.Title,
		"location":   schedule.Location,
		"comment":    schedule.Comment,
		"start_time": schedule.StartTime,
		"end_time":   schedule.EndTime,
	}

	result = dao.DB.Model(&existing).Updates(updateData)
	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, false, fmt.Errorf("日程更新失败")
	}

	// 返回更新后的数据（重新查询确保数据最新）
	dao.DB.First(&existing, schedule.Id)
	return &existing, true, nil
}

// DeleteSchedule 删除日程（根据ID和UserID验证权限）
func (dao *UserDao) DeleteSchedule(scheduleID, userID uint) (bool, error) {
	// 查询是否存在且属于该用户
	existing := model.Schedule{}
	result := dao.DB.Where("id = ? AND user_id = ?", scheduleID, userID).First(&existing)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("日程不存在或无权删除")
		}
		return false, result.Error
	}

	// 执行删除
	result = dao.DB.Delete(&existing)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("日程删除失败")
	}

	return true, nil
}
