package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nouhoum/casbin-go-example/api"
	"github.com/nouhoum/casbin-go-example/internal/service"
)

func CurrentUserID(c *gin.Context) int {
	user, _ := c.Get(service.UserIDKey)
	return user.(*api.AuthenticatedUser).ID
}
