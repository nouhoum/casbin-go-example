package service

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type CasbinConfig struct {
	CasbinModelFile string
	CasbinTableName string
}

func NewCasbinConfig(i *do.Injector) (*CasbinConfig, error) {
	return &CasbinConfig{
		CasbinModelFile: viper.GetString("CASBIN_CONFIG"),
		CasbinTableName: viper.GetString("CASBIN_RULES_TABLE"),
	}, nil
}

func NewCasbinEnforcer(i *do.Injector) (*casbin.Enforcer, error) {
	cfg := do.MustInvoke[*CasbinConfig](i)
	db := do.MustInvoke[*gorm.DB](i)
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &model.CasbinRule{}, cfg.CasbinTableName)
	if err != nil {
		return nil, fmt.Errorf("unable to create gorm adapter for postgres: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(cfg.CasbinModelFile, adapter)
	if err != nil {
		return nil, fmt.Errorf("unable to create casbin enforcer with adapter: %v", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("unable to load policy: %v", err)
	}

	return enforcer, nil
}
