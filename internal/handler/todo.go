package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nouhoum/casbin-go-example/api"
	"github.com/nouhoum/casbin-go-example/internal/service"
	"github.com/samber/do"
)

type Todo struct {
	service service.Todo
}

func NewTodo(i *do.Injector) (*Todo, error) {
	return &Todo{
		service: do.MustInvoke[service.Todo](i),
	}, nil
}

func (t *Todo) Create(c *gin.Context) {
	req := new(api.CreateOrUpdateTodoItemRequest)

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
	item, err := t.service.Create(ctx, req.Title, req.Description)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/todos/%d", item.ID))
	c.JSON(http.StatusCreated, item)
}

func (t *Todo) Get(c *gin.Context) {
	id := c.Param("id")

	item, err := t.service.Get(c.Request.Context(), id)
	if err == service.ErrTodoItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (t *Todo) Update(c *gin.Context) {
	req := new(api.CreateOrUpdateTodoItemRequest)

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

	id := c.Param("id")
	ctx := c.Request.Context()
	_, err = t.service.Get(ctx, id)
	if err == service.ErrTodoItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	item, err := t.service.Update(ctx, id, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (t *Todo) Delete(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	_, err := t.service.Get(ctx, id)
	if err == service.ErrTodoItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	item, err := t.service.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (t *Todo) List(c *gin.Context) {
	rawIsCompleted := c.Query("is_completed")
	var isCompleted *bool
	if rawIsCompleted != "" {
		result, err := strconv.ParseBool(rawIsCompleted)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid is_complete query parameter"})
			return
		}
		isCompleted = &result
	}

	page := parseInt(c.Query("page"), 1)
	size := parseInt(c.Query("size"), 100)

	items, total, err := t.service.List(
		service.WithPage(page),
		service.WithPageSize(size),
		service.WithIsCompleted(isCompleted),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"page":  page,
			"total": total,
		},
		"data": items,
	})
}

func (t *Todo) Complete(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	_, err := t.service.Get(ctx, id)
	if err == service.ErrTodoItemNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	item, err := t.service.UpdateCompleteness(ctx, id, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func parseInt(val string, defaultV int64) int64 {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultV
	}
	return i
}
