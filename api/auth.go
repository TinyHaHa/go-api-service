package api

import (
	"fmt"
	"github.com/google/wire"
	"github.com/opensourceai/go-api-service/internal/models"
	"github.com/opensourceai/go-api-service/internal/service"
	"github.com/opensourceai/go-api-service/middleware/jwt"
	"github.com/opensourceai/go-api-service/pkg/logging"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceai/go-api-service/pkg/app"
	"github.com/opensourceai/go-api-service/pkg/e"
	"github.com/opensourceai/go-api-service/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"` // 用户名
	Password string `valid:"Required; MaxSize(50)"` // 密码
}

type UserApi struct {
}

var ProviderAuth = wire.NewSet(NewAuthApi, service.ProviderUser)

func NewAuthApi(service2 service.UserService) (*UserApi, error) {
	userService = service2
	return &UserApi{}, nil
}

var userService service.UserService

//
func NewAuthRouter(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
	}
	group := router.Group("/auth")
	group.Use(jwt.JWT())
	{
		group.GET("/test", authTest)

	}

}

// @Summary 获取认证信息
// @Tags Auth
// @Produce  json
// @Param user body auth true "user"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/login [post]
func login(c *gin.Context) {
	user := auth{}
	appG := app.Gin{C: c}
	httpCode, errCode := app.BindAndValid(c, &user)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	userFound, isExist, err := userService.Login(models.User{Username: user.Username, Password: user.Password})
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}
	// 生成token
	token, err := util.GenerateToken(userFound.ID, userFound.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

// @Summary 添加用户
// @Tags Auth
// @Produce  json
// @Param user body models.User true "user"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/register [post]
func register(c *gin.Context) {
	user := models.User{}
	appG := app.Gin{C: c}
	httpCode, errCode := app.BindAndValid(c, &user)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if err := userService.Register(&user); err != nil {
		logging.Error(err)
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 用户密码修改
// @Tags Auth
// @Produce  json
// @Param user body models.User true "user"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth/updatePwd [post]
func updatePwd(c *gin.Context){
	user := models.User{}
	//var newPwd string
	//user.Password = newPwd
	appG := app.Gin{C:c}
	httpCode, errCode := app.BindAndValid(c, &user)
	if errCode != e.SUCCESS{
		appG.Response(httpCode, errCode, nil)
		return
	}
	if err := userService.UpdatePwd(user.Username, user.Password); err != nil{
		logging.Error(err)
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	//Name string `json:"name"  grom:"column:name;not null" valid:"Required;MaxSize(50)"`
	//// 密码
	//Password string `json:"password" grom:"column:password;not null" valid:"Required"`
	//// 描述
	//Description string `json:"description" grom:"column:description" valid:"MaxSize(200)"`
	//// 性别
	//Sex int `json:"sex" grom:"column:sex;not null" valid:"Min(1)"`
	//// 头像地址
	//AvatarSrc string `json:"avatar_src" grom:"column:avatar_src;not null"`
	//// 电子邮件
	//Email string `json:"email" grom:"column:email" valid:"Required;Email;MaxSize(100)"`
	//// 网站
	//WebSite string `json:"web_site" grom:"column:web_site" valid:"MaxSize(50)"`
	//// 公司
	//Company string `json:"company" grom:"column:company" valid:"MaxSize(50)"`
	//// 职位
	//Position string `json:"position" grom:"column:position" valid:"MaxSize(50)"`

}



// @Summary auth测试
// @Tags Auth
// @Produce  json
// @Param str query string true "string"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Security ApiKeyAuth
// @Router /auth/test [get]
func authTest(c *gin.Context) {
	query := c.Query("str")
	get, _ := c.Get("username")
	fmt.Println("get", get)
	str := fmt.Sprintf("%s", get)
	g := app.Gin{C: c}
	g.Response(http.StatusOK, e.SUCCESS, str+query)

}
