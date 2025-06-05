package usm

import (
	"errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nortoo/usm/model"
	"github.com/nortoo/usm/types"
	"gorm.io/gorm"
)

type (
	CasbinOptions struct {
		Store      *gorm.DB
		PolicyPath string
	}

	Options struct {
		Store         *gorm.DB
		CasbinOptions *CasbinOptions
	}
)

type Client struct {
	db        *gorm.DB
	casbinCli *casbin.Enforcer
}

var _ ClientService = (*Client)(nil)

type ClientService interface {
	CreateApplication(app *model.Application) error
	DeleteApplication(app *model.Application) error
	UpdateApplication(app *model.Application, cols ...string) error
	GetApplication(app *model.Application, cols ...interface{}) (*model.Application, error)
	ListApplications(q *types.QueryApplicationOptions) (ret []*model.Application, total int64, err error)

	CreateGroup(g *model.Group) error
	DeleteGroup(g *model.Group) error
	UpdateGroup(g *model.Group, cols ...string) error
	GetGroup(g *model.Group, cols ...interface{}) (*model.Group, error)
	ListGroups(q *types.QueryGroupOptions) (ret []*model.Group, total int64, err error)

	CreateMenu(m *model.Menu) error
	DeleteMenu(m *model.Menu) error
	UpdateMenu(m *model.Menu, cols ...string) error
	GetMenu(m *model.Menu, cols ...interface{}) (*model.Menu, error)
	ListMenus(q *types.QueryMenuOptions) (ret []*model.Menu, total int64, err error)

	CreateRole(r *model.Role) error
	DeleteRole(r *model.Role) error
	UpdateRole(r *model.Role, cols ...string) error
	GetRole(r *model.Role, cols ...interface{}) (*model.Role, error)
	ListRoles(q *types.QueryRoleOptions) (ret []*model.Role, total int64, err error)

	CreatePermission(p *model.Permission) error
	DeletePermission(p *model.Permission) error
	UpdatePermission(p *model.Permission, cols ...string) error
	GetPermission(p *model.Permission, cols ...interface{}) (*model.Permission, error)
	ListPermissions(q *types.QueryPermissionOptions) (ret []*model.Permission, total int64, err error)

	CreateUser(u *model.User) error
	DeleteUser(u *model.User) error
	UpdateUser(u *model.User, cols ...string) error
	GetUser(u *model.User, cols ...interface{}) (*model.User, error)
	ListUsers(q *types.QueryUserOptions) (ret []*model.User, total int64, err error)

	Authorize(role, tenant, resource, action string) (bool, error)
}

func newCasbinEnforcer(options *CasbinOptions) (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(options.Store)
	e, err := casbin.NewEnforcer(options.PolicyPath, a)
	if err != nil {
		return nil, err
	}

	return e, e.LoadPolicy()
}

// New creates a new USM client with the provided options,
// so that it can be used to manage users, roles, permissions, and other entities in your own applications.
func New(options *Options) (*Client, error) {
	if options.Store == nil {
		return nil, errors.New("store is required")
	}
	if options.CasbinOptions == nil {
		return nil, errors.New("casbin options is required")
	}
	if options.CasbinOptions.Store == nil {
		return nil, errors.New("casbin store is required")
	}
	if options.CasbinOptions.PolicyPath == "" {
		return nil, errors.New("casbin policy path is required")
	}

	err := model.RegisterModels(options.Store)
	if err != nil {
		return nil, err
	}

	e, err := newCasbinEnforcer(options.CasbinOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:        options.Store,
		casbinCli: e,
	}, err
}
