package service

import (
	"Graduation-Project/log"
	"Graduation-Project/model"
	"Graduation-Project/repo"
	"errors"
)

// 这里提供rpc服务端业务逻辑实现的方法 1、根据username查详情(登录验证+携带token获取信息 2、更改userinfo

type UserServiceHandler struct {
	MySQLHandler *repo.UserDao
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

func (handler *UserServiceHandler) GetLoginInfoByUsername(username string) (*model.User, bool, error) {
	user := model.User{}
	//查db
	dbUserInfo, ok, err := handler.MySQLHandler.GetUserInfoByUsername(username)
	if err != nil {
		log.Logger.Errorf("get db userinfo by username err: %v", err)
		return nil, ok, err
	}
	if ok {
		return dbUserInfo, ok, nil
	}
	//走到这就有问题了

	return &user, false, errors.New("user Invalid")
}
func (handler *UserServiceHandler) GetSchedulesByUserID(userId string) (*model.User, bool, error) {
	user := model.User{}
	//查db
	dbUserInfo, ok, err := handler.MySQLHandler.GetUserInfoByUsername(username)
	if err != nil {
		log.Logger.Errorf("get db userinfo by username err: %v", err)
		return nil, false, err
	}
	if ok {
		//这里还要写入redis
		return dbUserInfo, ok, nil
	}
	//走到这就有问题了

	return &user, false, errors.New("user Invalid")
}
