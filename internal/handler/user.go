package handler

import "github.com/gin-gonic/gin"

type User struct{}

func (u *User) Create(c *gin.Context)       {}
func (u *User) Authenticate(c *gin.Context) {}
