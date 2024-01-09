package perm

import (
	"context"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Initialize the model from a string.
const textModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act) || r.sub == "admin"
`

type CasbinUtil struct {
	_enforcer *casbin.Enforcer

	Logger *ylog.Logger
	DB     *gorm.DB
	Q      *query.Query
}

func NewCasbinUtil(
	logger *ylog.Logger,
	db *gorm.DB,
	q *query.Query) *CasbinUtil {
	return &CasbinUtil{
		Logger: logger,
		DB:     db,
		Q:      q,
	}
}

func (p *CasbinUtil) Enforcer() (*casbin.Enforcer, error) {
	if p._enforcer != nil {
		return p._enforcer, nil
	}

	ctx := context.Background()

	// model
	m, err := model.NewModelFromString(textModel)
	if err != nil {
		p.Logger.Errorf(ctx, "[error] model: %v", err)
		return nil, err
	}

	// adapter
	a, err := gormadapter.NewAdapterByDBUseTableName(p.DB, "t", "casbin")
	if err != nil {
		p.Logger.Errorf(ctx, "[error] adapter: %v", err)
		return nil, err
	}

	// enforcer
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		p.Logger.Errorf(ctx, "[error] enforcer: %v", err)
		return nil, err
	}
	p._enforcer = e

	return p._enforcer, nil
}
