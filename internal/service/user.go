package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nouhoum/casbin-go-example/internal/crypto"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
	Create(ctx context.Context, email, password, firstname, lastname string) (*model.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type userService struct {
	db *gorm.DB
}

func NewUser(i *do.Injector) (User, error) {
	return &userService{
		db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

func (svc *userService) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	fmt.Println("======> LOGIN=", login, "PASSWORD=", password)
	user := new(model.User)
	err := svc.db.Where(model.User{Email: login}).Take(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoSuchUser
		}
		return nil, err
	}

	err = crypto.Compare(user.Password, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidLoginOrPassword
		}
		return nil, err
	}

	return user, nil
}

func (svc *userService) Create(ctx context.Context, email string, password string, firstname string, lastname string) (*model.User, error) {
	email = strings.ToLower(email)
	isTaken, err := svc.IsEmailTaken(ctx, email)
	if err != nil {
		return nil, err
	}
	if isTaken {
		return nil, ErrEmailTaken
	}

	encryptedPwd, err := crypto.Encrypt(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
		Password:  encryptedPwd,
	}
	return user, svc.db.Create(user).Error
}

func (svc *userService) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	err := svc.db.Model(&model.User{}).
		Where("email = ?", strings.ToLower(email)).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, err
}
