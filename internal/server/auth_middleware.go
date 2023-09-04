package server

import (
	"fmt"
	"log"
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

		log.Println("======> ", subj, resource, act)
		ok, err := s.enforcer.Enforce(subj, resource, act)
		fmt.Println("ENFORCE RESULT=", ok)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !ok {
			fmt.Println("OOOPPS")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized action"})
			return
		}

		c.Next()
	}
}
