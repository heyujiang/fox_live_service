package middleware

import (
	"fmt"
	permissions "fox_live_service/internal/app/server/logic/permission"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"log"
	"sort"

	"github.com/gin-gonic/gin"
)

// ChkPer 根据权限判断用户是否有访问权限
func ChkPer(pl *permissions.PermissionsLogic, obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if obj == "" || act == "" {
			c.Next()
			return
		}

		sub := pl.SubFu(c)
		if sub == "" {
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("您没有执行该操作的权限哦！")))
			c.Abort()
			return
		}

		if ok, err := pl.Enforcer.Enforce(sub, obj, act); err != nil || !ok {
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("您没有执行该操作的权限哦！")))
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

// ChkRoles 根据角色判断用户是否有访问权限
func ChkRoles(pl *permissions.PermissionsLogic, requiredRoles []string, opts ...permissions.Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(requiredRoles) == 0 {
			c.Next()
			return
		}

		for _, opt := range opts {
			opt(pl)
		}
		sub := pl.SubFu(c)
		if sub == "" {
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("您没有执行该操作的权限哦！")))
			c.Abort()
			return
		}
		actualRoles, err := pl.Enforcer.GetRolesForUser(sub)
		if err != nil {
			log.Println("couldn't get roles of subject: ", err)
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("角色验证失败！")))
			c.Abort()
			return
		}

		sort.Strings(requiredRoles)
		sort.Strings(actualRoles)
		if pl.Opts.Logic == permissions.OR {
			for _, requiredRole := range requiredRoles {
				i := sort.SearchStrings(actualRoles, requiredRole)
				if i >= 0 && i < len(actualRoles) && actualRoles[i] == requiredRole {
					c.Next()
					return
				}
			}
			common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("您没有执行该操作的权限哦！")))
			c.Abort()
			return
		} else if pl.Opts.Logic == permissions.AND {
			for _, requiredRole := range requiredRoles {
				i := sort.SearchStrings(actualRoles, requiredRole)
				if i >= len(actualRoles) || actualRoles[i] != requiredRole {
					common.ResponseErr(c, errorx.NewErrorX(errorx.ErrCommon, fmt.Sprintf("您没有执行该操作的权限哦！")))
					c.Abort()
					return
				}
			}
			c.Next()
			return
		}
	}
}
