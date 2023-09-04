package service

import (
	"context"
	"errors"
	"time"

	"github.com/nouhoum/casbin-go-example/api"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type queryOptions struct {
	page           int64
	size           int64
	isCompleted    *bool
	skipPagination bool
}

type QueryOption func(*queryOptions)

func WithSkipPagination() QueryOption {
	return func(o *queryOptions) {
		o.skipPagination = true
	}
}

func WithPage(page int64) QueryOption {
	return func(o *queryOptions) {
		o.page = page
	}
}

func WithPageSize(size int64) QueryOption {
	return func(o *queryOptions) {
		o.size = size
	}
}

func WithIsCompleted(isCompleted *bool) QueryOption {
	return func(o *queryOptions) {
		o.isCompleted = isCompleted
	}
}

type Todo interface {
	Create(ctx context.Context, req api.CreateTodoItemRequest) (*model.TodoItem, error)
	Delete(ctx context.Context, id string) (*model.TodoItem, error)
	Get(ctx context.Context, id string) (*model.TodoItem, error)
	List(options ...QueryOption) ([]*model.TodoItem, int64, error)
	Update(ctx context.Context, id string, req api.UpdateTodoItemRequest) (*model.TodoItem, error)
}

type todoService struct {
	db     *gorm.DB
	policy Policy
}

func NewTodo(i *do.Injector) (Todo, error) {
	return &todoService{
		db:     do.MustInvoke[*gorm.DB](i),
		policy: do.MustInvoke[Policy](i),
	}, nil
}

func (svc *todoService) Get(ctx context.Context, id string) (*model.TodoItem, error) {
	item := new(model.TodoItem)

	if err := svc.db.First(item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (svc *todoService) List(options ...QueryOption) ([]*model.TodoItem, int64, error) {
	var items = make([]*model.TodoItem, 0)
	var total int64

	db := svc.applyOptions(options)

	if err := db.Find(&items).Error; err != nil {
		return []*model.TodoItem{}, 0, err
	}

	options = append(options, WithSkipPagination())
	db = svc.applyOptions(options)

	if err := db.Model(&items).Count(&total).Error; err != nil {
		return []*model.TodoItem{}, 0, err
	}

	return items, total, nil
}

func (svc *todoService) Delete(ctx context.Context, id string) (*model.TodoItem, error) {
	item, err := svc.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoItemNotFound
		}
		return nil, err
	}
	err = svc.db.Delete(item).Error
	if err != nil {
		return nil, err
	}

	err = svc.policy.OnTodoItemDeletion(ctx, *item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (svc *todoService) Update(ctx context.Context, id string, update api.UpdateTodoItemRequest) (*model.TodoItem, error) {
	item, err := svc.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoItemNotFound
		}
		return nil, err
	}

	item = updateItem(update, item)

	err = svc.db.Save(item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (svc *todoService) Create(ctx context.Context, req api.CreateTodoItemRequest) (*model.TodoItem, error) {
	item := &model.TodoItem{Title: req.Title, Description: req.Description, OwnerID: req.OwnerID}
	err := svc.db.Create(item).Error
	if err != nil {
		return nil, err
	}
	err = svc.policy.OnTodoItemCreation(ctx, *item)
	if err != nil {
		return nil, err
	}

	return item, err
}

func updateItem(update api.UpdateTodoItemRequest, item *model.TodoItem) *model.TodoItem {
	if update.Title != nil {
		item.Title = *update.Title
	}

	if update.Description != nil {
		item.Description = *update.Description
	}

	if update.Complete != nil {
		if *update.Complete {
			now := time.Now()
			item.CompletedAt = &now
		} else {
			item.CompletedAt = nil
		}
	}
	return item
}

func (svc *todoService) applyOptions(options []QueryOption) *gorm.DB {
	db := svc.db

	opts := &queryOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if !opts.skipPagination {
		if opts.page > 0 {
			db = db.Offset(int((opts.page - 1) * opts.size))
		}

		if opts.size > 0 {
			db = db.Limit(int(opts.size))
		}
	}

	if opts.isCompleted != nil {
		if *opts.isCompleted {
			db = db.Where("completed_at is not null")
		} else {
			db = db.Where("completed_at is null")
		}
	}

	return db
}
