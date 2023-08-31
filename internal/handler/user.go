package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nouhoum/casbin-go-example/api"
	"github.com/nouhoum/casbin-go-example/internal/service"
	"github.com/samber/do"
)

type User struct {
	service service.User
}

func NewUser(i *do.Injector) (*User, error) {
	return &User{
		service: do.MustInvoke[service.User](i),
	}, nil
}

func (u *User) Create(c *gin.Context) {
	req := new(api.CreateUserRequest)

	err := c.ShouldBindBodyWith(req, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := u.service.Create(ctx, req.Email, req.Password, req.Firstname, req.Lastname)

	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/users/%d", user.ID))
	c.JSON(http.StatusCreated, user)
}
