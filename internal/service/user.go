package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/nouhoum/casbin-go-example/api"
	"github.com/nouhoum/casbin-go-example/internal/crypto"
	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
	Create(ctx context.Context, req api.CreateUserRequest) (*model.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	InitUsers() error // InitUsers for test purposes only
}

type userService struct {
	db     *gorm.DB
	policy Policy
}

func NewUser(i *do.Injector) (User, error) {
	return &userService{
		db:     do.MustInvoke[*gorm.DB](i),
		policy: do.MustInvoke[Policy](i),
	}, nil
}

func (service *userService) InitUsers() error {
	var count int64
	err := service.db.Model(&model.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		log.Println("intitializing users...")

		users := []api.CreateUserRequest{
			{Email: "john.doe@todo.com", Password: "secret", Firstname: "John", Lastname: "Doe", RoleID: &Admin.ID},
			{Email: "jane.doe@todo.com", Password: "secret", Firstname: "Jane", Lastname: "Doe", RoleID: &SuperAdmin.ID},
			{Email: "dupont@todo.com", Password: "secret", Firstname: "Jean", Lastname: "Dupont", RoleID: &Member.ID},
		}

		for _, user := range users {
			_, err = service.Create(context.Background(), user)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *userService) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
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

func (svc *userService) Create(ctx context.Context, req api.CreateUserRequest) (*model.User, error) {
	email := strings.ToLower(req.Email)
	isTaken, err := svc.IsEmailTaken(ctx, email)
	if err != nil {
		return nil, err
	}
	if isTaken {
		return nil, ErrEmailTaken
	}

	encryptedPwd, err := crypto.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  encryptedPwd,
		RoleID:    *req.RoleID,
	}

	err = svc.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	err = svc.db.Preload("Role").First(user).Error
	if err != nil {
		return nil, err
	}
	err = svc.policy.OnUserCreation(ctx, *user)
	if err != nil {
		return nil, err
	}

	return user, nil
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
