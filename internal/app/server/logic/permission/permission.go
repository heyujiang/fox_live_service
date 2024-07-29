package permissions

import (
	"fmt"
	"fox_live_service/config/global"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"log"
)

var Permission *PermissionsLogic

func init() {
	policyAdapter, err := gormadapter.NewAdapter(
		"mysql",
		global.Config.GetString("Db.Mysql.DSN"),
		true)
	if err != nil {
		log.Fatalf(err.Error())
	}

	Permission, err = newPermissionLogic(func(c *gin.Context) string {
		return c.GetString("username")
	}, global.Config.GetString("Casbin.ModelFile"), policyAdapter)

	if err != nil {
		log.Fatalf(err.Error())
	}
}

type Logic int

const (
	OR Logic = iota
	AND
)

type Option func(r *PermissionsLogic)

type options struct {
	Logic Logic
}

func WithLogic(logic Logic) Option {
	return func(r *PermissionsLogic) {
		r.Opts.Logic = logic
	}
}

type (
	PermissionsLogic struct {
		Enforcer *casbin.Enforcer
		SubFu    SubFn
		Opts     *options
	}

	SubFn func(c *gin.Context) string
)

func newPermissionLogic(fn SubFn, modelFile string, policyAdapter interface{}) (*PermissionsLogic, error) {
	e, err := newCasbinEnforcer(modelFile, policyAdapter)
	if err != nil {
		return nil, err
	}
	if fn == nil {
		slog.Error("casbin permission middleware get subject func cannot be nil")
		return nil, fmt.Errorf("fn cannot be nil")
	}

	return &PermissionsLogic{
		Enforcer: e,
		SubFu:    fn,
		Opts: &options{
			Logic: OR,
		},
	}, nil
}

func newCasbinEnforcer(modelFile string, policyAdapter interface{}) (*casbin.Enforcer, error) {
	enforcer, err := casbin.NewEnforcer(modelFile, policyAdapter)
	if err != nil {
		return nil, fmt.Errorf("casbin enforce fail,err:%v", err)
	}
	return enforcer, nil
}

func (pl *PermissionsLogic) DeleteRoleForUser(user string, role string, domain ...string) error {
	isSuccess, err := pl.Enforcer.DeleteRoleForUser(user, role, domain...)
	if err != nil {
		return fmt.Errorf("casbin del role:[%s] for user:[%s] fail,err:%v", user, role, err)
	}
	if !isSuccess {
		return fmt.Errorf("casbin del role:[%s] for user:[%s] fail ", user, role)
	}

	return nil
}

func (pl *PermissionsLogic) AddRoleForUser(user string, role string, domain ...string) error {
	isSuccess, err := pl.Enforcer.AddRoleForUser(user, role, domain...)
	if err != nil {
		return fmt.Errorf("casbin add role:[%s] for user:[%s] fail,err:%v", user, role, err)
	}
	if !isSuccess {
		return fmt.Errorf("casbin add role:[%s] for user:[%s] fail ", user, role)
	}
	return nil
}

func (pl *PermissionsLogic) GetUsersForRole(name string, domain ...string) ([]string, error) {
	users, err := pl.Enforcer.GetUsersForRole(name, domain...)
	if err != nil {
		return []string{}, fmt.Errorf("casbin gets the users that has a role:[%s] , err : %v", name, err)
	}

	return users, nil
}

func (pl *PermissionsLogic) GetRolesForUser(name string, domain ...string) ([]string, error) {
	roles, err := pl.Enforcer.GetRolesForUser(name, domain...)
	if err != nil {
		return []string{}, fmt.Errorf("casbin gets the roles that has a user:[%s] , err : %v", name, err)
	}

	return roles, nil
}
