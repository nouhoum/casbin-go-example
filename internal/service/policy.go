package service

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
)

type Action = string

const (
	Write  Action = "write"
	Delete Action = "delete"
	Read   Action = "read"
)

var Actions = []Action{Read, Write, Delete}

type Policy interface {
	OnUserDeletion(ctx context.Context, item model.User) error
	OnUserRoleAssignment(ctx context.Context, userID, roleID int) error
	OnUserRoleDeletion(ctx context.Context, userID, roleID int) error
	OnUserCreation(ctx context.Context, item model.User) error
	OnTodoItemCreation(ctx context.Context, item model.TodoItem) error
	OnTodoItemDeletion(ctx context.Context, item model.TodoItem) error
}

type policy struct {
	enforcer *casbin.Enforcer
}

func NewPolicy(i *do.Injector) (Policy, error) {
	return &policy{
		enforcer: do.MustInvoke[*casbin.Enforcer](i),
	}, nil
}

func (p *policy) OnTodoItemCreation(ctx context.Context, item model.TodoItem) error {
	rules := [][]string{}
	subj := fmt.Sprint(item.OwnerID)
	obj := fmt.Sprintf("todos.%d", item.ID)
	for _, act := range Actions {
		rules = append(rules, []string{subj, obj, act})
	}

	_, err := p.enforcer.AddPolicies(rules)
	if err != nil {
		return err
	}

	return nil
}

func (p *policy) OnTodoItemDeletion(ctx context.Context, item model.TodoItem) error {
	rules := [][]string{}
	subj := fmt.Sprint(item.OwnerID)
	obj := fmt.Sprint(item.ID)
	for _, act := range Actions {
		rules = append(rules, []string{subj, obj, act})
	}

	_, err := p.enforcer.RemovePolicies(rules)
	if err != nil {
		return err
	}

	return nil
}

func (p *policy) OnUserCreation(ctx context.Context, user model.User) error {
	_, err := p.enforcer.AddGroupingPolicy(fmt.Sprint(user.ID), user.Role.Slug)
	if err != nil {
		return err
	}

	return nil
}

func (p *policy) OnUserDeletion(ctx context.Context, user model.User) error {
	panic("unimplemented")
}

// OnUserRoleAssignment implements Policy.
func (p *policy) OnUserRoleAssignment(ctx context.Context, userID int, roleID int) error {
	panic("unimplemented")
}

// OnUserRoleDeletion implements Policy.
func (p *policy) OnUserRoleDeletion(ctx context.Context, userID int, roleID int) error {
	panic("unimplemented")
}
