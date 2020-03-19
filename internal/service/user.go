package service

import (
	"errors"
	"github.com/google/wire"
	"github.com/opensourceai/go-api-service/internal/dao"
	"github.com/opensourceai/go-api-service/internal/dao/mysql"
	"github.com/opensourceai/go-api-service/internal/models"
	"github.com/opensourceai/go-api-service/pkg/util"
)

type UserService interface {
	Register(user *models.User) error
	Login(user models.User) (*models.User, bool, error)
	//MsgEdit(user models.User) error
	UpdatePwd(username string, s string) error
}

type userService struct {
	dao.UserDao
}

var ProviderUser = wire.NewSet(NewUserService, mysql.NewUserDao)

func NewUserService(dao2 dao.UserDao) (UserService, error) {
	return &userService{dao2}, nil

}
func (service userService) Register(user *models.User) error {
	// 加密密码
	user.Password = util.EncodeMD5(user.Password)
	return service.DaoAdd(user)
}

func (service userService) Login(user models.User) (*models.User, bool, error) {
	err, u := service.DaoGetUserByUsername(user.Username)
	if err != nil {
		return nil, false, errors.New("登录失败")
	}
	// 匹配密码
	md5Password := util.EncodeMD5(user.Password)
	if u.Password == md5Password {
		return &u, true, nil
	}
	return nil, false, errors.New("登录失败")
}

func (service userService) UpdatePwd(username string, s string) error{
	//通过用户名从数据库获取用户对象
	_, u := service.DaoGetUserByUsername(username)
	//修改密码
	u.Password = s
	//调用修改用户信息方法将对象重新写入数据库，有错误就返回错误
	return service.DaoEdit(&u)
}

//func (service userService) MsgEdit(user models.User) error{
//	user
//	return service.DaoEdit()
//}