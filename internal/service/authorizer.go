package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/samber/do"
)

type Authorizer interface {
	IsAuthorized(subj, resource, action string) (bool, error)
}

type authorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(i *do.Injector) (Authorizer, error) {
	return &authorizer{
		enforcer: do.MustInvoke[*casbin.Enforcer](i),
	}, nil
}

func (a authorizer) IsAuthorized(subj, resource, action string) (bool, error) {
	return a.enforcer.Enforce(subj, resource, action)
}
