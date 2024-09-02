package handler

import (
	"github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// userRegister 用户注册
// @Tags  user
// @Summary 用户注册
// @Param req body dto.UserRegisterReq true "管理员注册提交参数"
// @Router /user/public/register [post]
// @Success 200 {object} conf.Response
// @Failure 400 {object} string "参数错误"
// @Failure 500 {object} string "内部错误"
// @Produce json
// @Accept json
func userRegister(c *gin.Context) {
	r := new(dto.UserRegisterReq)

	if err := c.Bind(r); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}
	userRegisterResp, err := svc.UserRegister(c.Request.Context(), r.Wallet)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}
	if userRegisterResp.UserExists == true {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "fail", Data: "钱包重复"})
		return
	}
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: userRegisterResp.User})
}

// userLogin 用户登录
// @Tags  user
// @Summary 用户登录
// @Param req body dto.UserLoginReq true "用户登录提交参数"
// @Router /user/public/login [post]
// @Success 200 {object} conf.Response
// @Failure 400 {object} string "参数错误"
// @Failure 500 {object} string "内部错误"
// @Produce json
// @Accept json
func userLogin(c *gin.Context) {
	r := new(dto.UserLoginReq)

	if err := c.Bind(r); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}

	res, err := svc.UserLogin(c.Request.Context(), r.Wallet)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}
	if res.Token == "用户不存在" {
		c.JSON(http.StatusBadRequest, conf.Response{Code: 400, Msg: "fail", Data: res.Token})
		return
	}
	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: res.Token})
}

// userQuery 用户查询
// @Tags  user
// @Summary 用户查询
// @Param wallet query string true "钱包地址"
// @Router /user/protected/query [get]
// @Success 200 {object} conf.Response
// @Failure 400 {object} conf.ResponseError
// @Failure 500 {object} string "内部错误"
func userQuery(c *gin.Context) {
	r := new(dto.UserQueryReq)

	if err := c.Bind(r); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}
	query, err := svc.UserQuery(c.Request.Context(), r.Wallet)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, conf.Response{Code: 500, Msg: "err", Data: err.Error()})
		return
	}

	c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: query})
}
