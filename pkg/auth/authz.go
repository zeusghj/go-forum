package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Authz 定义了一个授权器，提供授权功能.
type Authz struct {
	*casbin.SyncedEnforcer
}

// NewAuthz 创建一个使用 casbin 完成授权的授权器.
func NewAuthz(db *gorm.DB) (*Authz, error) {

	// Initialize a Gorm adapter and use it in a Casbin enforcer
	adapter, err := gormadapter.NewAdapterByDB(db) // 用 GORM 存储策略
	if err != nil {
		return nil, err
	}

	// Initialize the enforcer.
	enforcer, err := casbin.NewSyncedEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// load the policy from DB.
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	// refresh every 30 seconds
	enforcer.StartAutoLoadPolicy(30 * time.Second)

	a := &Authz{enforcer}

	return a, nil

}

// Authorize 用来进行授权.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
