package service

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"gorm.io/gorm"
)

var (
	Admin      = model.Role{ID: 1, Slug: "admin"}
	SuperAdmin = model.Role{ID: 2, Slug: "super-admin"}
	Member     = model.Role{ID: 3, Slug: "member"}
)

var initialRoles = []model.Role{Admin, SuperAdmin, Member}
var AdminRoles = []model.Role{Admin, SuperAdmin}
var AdminRolesMappings = map[string][]Action{
	Admin.Slug:      {Read, Write},
	SuperAdmin.Slug: {Read, Write, Delete},
}

type Role interface {
	InitRoles() error
}

type role struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

func NewRole(i *do.Injector) (Role, error) {
	return &role{
		db:       do.MustInvoke[*gorm.DB](i),
		enforcer: do.MustInvoke[*casbin.Enforcer](i),
	}, nil
}

func (service *role) InitRoles() error {
	var count int64
	err := service.db.Model(&model.Role{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		log.Println("intitializing roles...")
		err = service.db.CreateInBatches(initialRoles, len(initialRoles)).Error
		if err != nil {
			return err
		}
	}

	for s, as := range AdminRolesMappings {
		for _, a := range as {
			service.enforcer.AddPolicy(s, "*", a)
		}
	}

	return nil
}
