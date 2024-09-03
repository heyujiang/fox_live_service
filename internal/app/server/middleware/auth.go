package middleware

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"fox_live_service/pkg/util/jwttokenx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			slog.Error("token is empty")
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "Authorization 不能为空"))
			c.Abort()
			return
		}

		authList := strings.Split(tokenString, " ")
		if len(authList) != 2 {
			slog.Error("token 格式错误", "token", tokenString)
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "Authorization 格式错误"))
			c.Abort()
			return
		}

		if authList[0] != "Bearer" {
			slog.Error("token 请以Bearer开头", "token", tokenString)
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "Authorization 格式错误"))
			c.Abort()
			return
		}

		if authList[1] == "null" {
			slog.Error("token 不合法", "token", tokenString)
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "Authorization 格式错误"))
			c.Abort()
			return
		}

		userId, username, expireIn, refreshExpireIn, err := jwttokenx.ParseValuePasswordFromToken(authList[1], global.Config.GetString("AccessToken.JwtTokenKey"))

		if err != nil {
			if errors.Is(err, jwttokenx.ErrJwtTokenExpired) {
				slog.Error("token已过期, 阶段1", "token", tokenString, "expireIn", expireIn, "refreshExpireIn", refreshExpireIn)
				common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "需要重新登录"))
				c.Abort()
				return
			} else {
				slog.Error("token已过期, 阶段2", "token", tokenString, "expireIn", expireIn, "refreshExpireIn", refreshExpireIn)
				common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "无效的token"))
				c.Abort()
			}
			c.Abort()
			return
		}

		//验证是否过期
		//判断刷新token是否过期 : 如果过期
		//过期重新登录
		//未过期 刷新token

		user, err := model.UserModel.Find(userId)
		if err != nil {
			slog.Error("find user error", "userId", userId, "err", err)
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrNoLogin, "验证用户失败"))
			c.Abort()
		}
		isSuper := false
		roleIds := strings.Split(user.RoleIds, ",")
		for _, roleId := range roleIds {
			if cast.ToInt(roleId) == model.SuperManagerRoleId {
				isSuper = true
			}
		}

		c.Set("uid", userId)
		c.Set("username", username)
		c.Set("isSuper", isSuper)

		c.Next()
	}
}
