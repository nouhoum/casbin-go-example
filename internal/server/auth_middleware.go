package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nouhoum/casbin-go-example/internal/handler"
)

func (s *Server) AccessControl(obj, act, objIDName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj := fmt.Sprint(handler.CurrentUserID(c))

		var objID string
		resource := obj
		if objIDName != "" {
			objID = c.Param(objIDName)
			resource = fmt.Sprintf("%s.%s", resource, objID)
		}

		ok, err := s.enforcer.Enforce(subj, resource, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized action"})
			return
		}

		c.Next()
	}
}
