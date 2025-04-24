package service

import (
	"Graduation-Project/config"
	"Graduation-Project/log"
	"Graduation-Project/model"
	"Graduation-Project/repo"
	"Graduation-Project/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// 这里提供rpc服务端业务逻辑实现的方法 1、根据username查详情(登录验证+携带token获取信息 2、更改userinfo

type UserServiceHandler struct {
	MySQLHandler *repo.UserDao
}

type ScheduleService interface {
	TextToSchedule(text string) (utils.ScheduleJson, error)
}

type DeepSeekScheduleService struct {
	apiClient   *http.Client
	apiEndpoint string
	apiKey      string
}

func NewScheduleService() (ScheduleService, error) {
	return &DeepSeekScheduleService{
		apiClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiEndpoint: "https://api.deepseek.com/v1/schedule/parse",
		apiKey:      config.GetConfig().DeepSeekAPIKey,
	}, nil
}

func (s *DeepSeekScheduleService) TextToSchedule(text string) (utils.ScheduleJson, error) {
	schedule, err := utils.TextToJson(text)
	if err != nil {
		return schedule, err
	}
	return schedule, nil
}
func NewUserServiceHandler() (handler *UserServiceHandler, err error) {
	handler = &UserServiceHandler{}
	dbHandler, err := repo.InitDB()
	if err != nil {
		log.Logger.Errorf("db init err: %v", err)
		return handler, err
	}
	handler.MySQLHandler = &repo.UserDao{DB: dbHandler}

	return handler, nil
}
func (handler *UserServiceHandler) Close() {
	sqlDB, err := handler.MySQLHandler.DB.DB()
	if err != nil {

	} else {
		// 关闭数据库连接
		err = sqlDB.Close()
		if err != nil {
			log.Logger.Errorf("close db err: %v", err)
		}
	}

}

// Login 用户登录逻辑
func (handler *UserServiceHandler) Login(username, password string) (uint, error) {
	// 参数基础校验
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return 0, errors.New("用户名和密码不能为空")
	}

	// 查询用户信息
	user, exists, err := handler.MySQLHandler.GetUserInfoByUsername(username)
	if err != nil {
		return 0, fmt.Errorf("登录失败：数据库查询错误 (%v)", err)
	}
	if !exists {
		return 0, errors.New("用户名不存在")
	}

	// 使用数据库中的盐值加密输入密码
	hashedPassword := utils.Md5BySalt(password, user.Salt) // user.Salt = 用户名

	// 对比密码哈希值
	if hashedPassword != user.Password {
		return 0, errors.New("密码错误")
	}

	// 登录成功，返回用户ID
	return user.Id, nil
}
func (handler *UserServiceHandler) GetSchedulesByUserID(userId uint) ([]model.Schedule, bool, error) {
	schedule := []model.Schedule{}
	//查db
	dbUserInfo, ok, err := handler.MySQLHandler.GetScheduleInfoByUserId(userId)
	if err != nil {
		log.Logger.Errorf("get db userinfo by username err: %v", err)
		return nil, false, err
	}
	if ok {
		//这里还要写入redis
		return dbUserInfo, ok, nil
	}
	//走到这就有问题了

	return schedule, false, nil
}
func (handler *UserServiceHandler) CreateUser(username, password string) error {
	// 参数基础校验
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("用户名和密码不能为空")
	}

	// 生成带盐哈希密码（盐=用户名）
	hashedPassword := utils.Md5BySalt(password, username)

	// 构造用户对象
	newUser := &model.User{
		Username: username,
		Password: hashedPassword,
		Salt:     username, // 根据Md5BySalt逻辑，盐值=用户名
	}

	// 调用DAO层插入数据库
	_, _, err := handler.MySQLHandler.CreateUser(newUser)
	if err != nil {
		// 直接透传DAO层封装好的错误（如"用户名xxx已存在"）
		return err
	}

	return nil
}

// ------------------------ 创建日程 ------------------------
func (handler *UserServiceHandler) CreateSchedule(schedule *model.Schedule) error {
	// 参数基础校验
	if strings.TrimSpace(schedule.Title) == "" {
		return errors.New("日程标题不能为空")
	}
	if schedule.StartTime == "" || schedule.EndTime == "" {
		return errors.New("必须指定起止时间")
	}
	if schedule.UserID == 0 {
		return errors.New("无效的用户标识")
	}

	// 调用DAO层
	_, _, err := handler.MySQLHandler.CreateSchedule(schedule)
	if err != nil {
		// 直接透传DAO层错误（如"用户ID不存在"、"日程冲突"等）
		return err
	}
	return nil
}

// ------------------------ 更新日程 ------------------------
func (handler *UserServiceHandler) UpdateSchedule(requesterUserID uint, schedule *model.Schedule) error {
	// 参数基础校验
	if schedule.Id == 0 {
		return errors.New("无效的日程标识")
	}
	if strings.TrimSpace(schedule.Title) == "" {
		return errors.New("日程标题不能为空")
	}
	if schedule.StartTime == "" || schedule.EndTime == "" {
		return errors.New("必须指定起止时间")
	}

	// 权限校验：确保操作者与日程所属用户一致
	if schedule.UserID != requesterUserID {
		return errors.New("无权修改他人日程")
	}

	// 调用DAO层
	_, _, err := handler.MySQLHandler.UpdateSchedule(schedule)
	if err != nil {
		return fmt.Errorf("更新失败: %v", err)
	}
	return nil
}

// ------------------------ 删除日程 ------------------------
func (handler *UserServiceHandler) DeleteSchedule(requesterUserID uint, scheduleID uint) error {
	// 参数基础校验
	if scheduleID == 0 {
		return errors.New("无效的日程标识")
	}

	// 调用DAO层（包含权限验证）
	success, err := handler.MySQLHandler.DeleteSchedule(scheduleID, requesterUserID)
	if err != nil {
		if strings.Contains(err.Error(), "无权删除") {
			return errors.New("无权操作该日程")
		}
		return fmt.Errorf("删除失败: %v", err)
	}
	if !success {
		return errors.New("日程删除失败")
	}
	return nil
}
