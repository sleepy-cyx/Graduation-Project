package service

import (
	"Graduation-Project/log"
	"Graduation-Project/model"
	"Graduation-Project/repo"
	"context"
	"errors"
	"time"
)

// 这里提供rpc服务端业务逻辑实现的方法 1、根据username查详情(登录验证+携带token获取信息 2、更改userinfo

type UserServiceHandler struct {
	MySQLHandler *repo.UserDao
	RedisHandler *repo.RedisDB
}

func NewUserServiceHandler() (handler *UserServiceHandler, err error) {
	handler = &UserServiceHandler{}
	handler.RedisHandler, err = repo.InitRedisDefault()
	if err != nil {
		println(err.Error())
		log.Logger.Errorf("redis init err: %v", err)
		return handler, err
	}
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

	if handler.RedisHandler != nil {
		err = handler.RedisHandler.Close()
		if err != nil {
			log.Logger.Errorf("close redis err: %v", err)
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
		//这里还要写入redis
		//err = handler.RedisHandler.SetByUsername(username, dbUserInfo)
		//if err != nil {
		//	Logger.Errorf("set redis userinfo by username err: %v", err)
		//}
		return dbUserInfo, ok, nil

	}
	//走到这就有问题了

	return &user, false, errors.New("user Invalid")
}
func (handler *UserServiceHandler) GetUserInfoByUsername(username string) (*model.User, bool, error) {
	user := model.User{}
	//先查redis，再查db
	ok, cacheInfo, err := handler.RedisHandler.GetByUsername(username, ctx)
	if err != nil {
		log.Logger.Errorf("get redis userinfo by username err: %v", err)
		return nil, false, err
	}
	if ok {
		return &cacheInfo, ok, nil
	}
	//查db
	dbUserInfo, ok, err := handler.MySQLHandler.GetUserInfoByUsername(username)
	if err != nil {
		log.Logger.Errorf("get db userinfo by username err: %v", err)
		return nil, false, err
	}
	if ok {
		//这里还要写入redis
		err = handler.RedisHandler.SetByUsername(username, dbUserInfo, ctx)
		if err != nil {
			log.Logger.Errorf("set redis userinfo by username err: %v", err)
		}
		return dbUserInfo, ok, nil

	}
	//走到这就有问题了

	return &user, false, errors.New("user Invalid")
}

func (handler *UserServiceHandler) AlterUserInfoByUsername(username string) (bool, error) {
	loginTimeParsed, err := time.Parse(time.RFC3339, req.LoginTime)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	err = handler.RedisHandler.DeleteInfo(ctx, username)
	if err != nil {
		return false, err
	}
	exist, err := handler.MySQLHandler.AlterUserInfo(username, req.Nickname, req.PictureUrl, req.LoginIp, loginTimeParsed)
	if err != nil {
		log.Logger.Errorf("alter userinfo by username err: %v", err)
		return false, err
	}
	if !exist {
		return false, nil
	}

	return true, nil
}
