package v1

import (
	"douyin_service/global"
	"douyin_service/internal/service"
	"douyin_service/pkg/app"
	"douyin_service/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)
// Login 登录
func (u User) Login(c *gin.Context) {
	param := service.LoginRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c.Request.Context())
	userId, flag, err := svc.Login(&param)
	res := &service.LoginResponse{
		UserID: userId,
		Token:  "",
	}
	res.StatusCode = errcode.ErrorLoginFail.Code()
	res.StatusMsg = errcode.ErrorLoginFail.Msg()
	if err != nil {
		global.Logger.Errorf("svc.Login err: %v", err)
		response.ToResponse(res)
		return
	}

	if !flag {
		global.Logger.Error("用户名/密码错误")
		response.ToResponse(res)
		return
	}
	idStr := strconv.Itoa(int(userId))
	token, err := app.GenerateToken(global.JWTSetting.Key, global.JWTSetting.Secret, idStr)
	if err != nil {
		global.Logger.Errorf("app.GenerateToken err: %v", err)
		response.ToErrorResponse(errcode.ErrorRegisterFail)
		return
	}
	res = &service.LoginResponse{
		UserID: userId,
		Token:  token,
	}
	res.StatusCode = 0
	res.StatusMsg = "登录成功"
	response.ToResponse(res)
	//return	//多余的return
}
