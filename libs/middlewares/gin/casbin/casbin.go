package CCasbin

import (
	"errors"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	CGinResponse "github.com/zander-84/go-components/libs/middlewares/gin/response"
	"log"
	"strings"
)

type Rbac struct {
	mysql *gorm.DB
	Obj   *casbin.Enforcer
}

func New() *Rbac {
	return new(Rbac)
}

func (this *Rbac) Init(model string, db *gorm.DB) *Rbac {
	this.mysql = db
	if model == "" {
		model =
			`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "root"
`
	}

	m := casbin.NewModel(model)
	this.Obj = casbin.NewEnforcer(m, gormadapter.NewAdapterByDB(db))
	if err := this.Obj.LoadPolicy(); err != nil {
		log.Fatalln("rbac LoadPolicy error")
	}
	return this
}
func (this *Rbac) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, ok := c.Get("casbin_role"); !ok {
			CGinResponse.Resp.Forbidden(c)
		} else {
			roleStr := role.(string)
			if roleStr != "root" {
				if ok, err := this.Obj.EnforceSafe(roleStr, c.FullPath(), strings.ToUpper(c.Request.Method)); !ok || err != nil {
					CGinResponse.Resp.Forbidden(c)
				}
			}
		}

		if CGinResponse.HasResp(c) {
			CGinResponse.HandleResponse(c, false)
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

type Policy struct {
	Path   string `json:"path" form:"path" validate:"required" comment:"请求路径"`
	Method string `json:"method"  form:"method" validate:"required" comment:"请求方法"`
}

// 添加页面对应的接口
func (this *Rbac) AddPagePolicies(page string, policies []Policy) error {
	if err := this.IsRoot(page); err != nil {
		return err
	}
	for _, p := range policies {
		_, err := this.Obj.AddPolicySafe(page, p.Path, strings.ToUpper(p.Method))
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除页面对应的接口
func (this *Rbac) DeletePagePolices(page string, policies []Policy) error {
	if err := this.IsRoot(page); err != nil {
		return err
	}

	for _, p := range policies {
		_, err := this.Obj.RemovePolicySafe(page, p.Path, strings.ToUpper(p.Method))
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取页面对应的接口
func (this *Rbac) GetPagePolices(page string) [][]string {
	return this.Obj.GetPermissionsForUser(page)
}

// 获取所有页面对应的接口
func (this *Rbac) GetPagesPolices() map[string][][]string {
	data := this.Obj.GetPolicy()
	pages := make(map[string][][]string, 0)

	for _, val := range data {
		pages[val[0]] = append(pages[val[0]], val)
	}
	return pages
}

// 删除页面
func (this *Rbac) DeletePages(pages []string) {
	for _, page := range pages {
		this.Obj.DeleteRole(page)
	}
}

// 添加角色页面
func (this *Rbac) AddRolePages(role string, pages []string) {
	for _, page := range pages {
		this.Obj.AddRoleForUser(role, page)
	}
}

// 删除角色页面
func (this *Rbac) DeleteRolePages(role string, pages []string) {
	for _, page := range pages {
		this.Obj.DeleteRoleForUser(role, page)
	}
}

// 获取角色页面
func (this *Rbac) GetRolePages(role string) (pages []string, err error) {
	return this.Obj.GetRolesForUser(role)
}

// 获取所有角色 用户
func (this *Rbac) GetRoles() (pages []string, err error) {
	var data []string
	err = this.mysql.Table("casbin_rule").Where("p_type = ?", "g").Group("v0").Pluck("v0", &data).Error
	return data, err
}

func (this *Rbac) IsRoot(name string) error {
	if name == "root" {
		return errors.New("Name Can't  root")
	}
	return nil
}
