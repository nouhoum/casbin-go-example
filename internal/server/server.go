package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nouhoum/casbin-go-example/internal/handler"
	"github.com/rs/cors"
	"github.com/samber/do"
	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	AllowedOrigins []string
}

func NewConfig(i *do.Injector) (*Config, error) {
	viper.AutomaticEnv()
	port := viper.GetInt("PORT")

	origins := strings.Split(viper.GetString("ALLOWED_ORIGINS"), ",")
	return &Config{
		Port:           fmt.Sprint(port),
		AllowedOrigins: origins,
	}, nil
}

type Server struct {
	Cfg        *Config
	HTTPServer *http.Server
	engine     *gin.Engine

	user *handler.User
	todo *handler.Todo
}

func New(i *do.Injector) (*Server, error) {
	engine := do.MustInvoke[*gin.Engine](i)
	gin.SetMode(gin.DebugMode)

	cfg := do.MustInvoke[*Config](i)

	opts := cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowedHeaders: []string{
			"authorization",
			"accept",
			"content-type",
			"Origin", "Accept", "Content-Type", "X-Requested-With",
		},
	}

	s := Server{
		Cfg: cfg,
		HTTPServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: cors.New(opts).Handler(engine),
		},
		engine: engine,
		todo:   do.MustInvoke[*handler.Todo](i),
	}

	s.addRoutes()
	return &s, nil
}

func (s *Server) Run() error {
	log.Printf("server listening on port %s", s.Cfg.Port)
	return s.HTTPServer.ListenAndServe()
}

func NewEngine(i *do.Injector) (*gin.Engine, error) {
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Up"})
	})

	return r, nil
}

func (s *Server) addRoutes() {
	api := s.engine.Group("/api")

	todos := api.Group("/todos")
	{
		todos.GET("/:id", s.todo.Get)
		todos.GET("", s.todo.List)
		todos.POST("", s.todo.Create)
		todos.POST("/:id", s.todo.Update)
		todos.POST("/:id/complete", s.todo.Complete)
		todos.DELETE("/:id", s.todo.Delete)
	}

	users := api.Group("/users")
	{
		users.POST("", s.user.Create)
		users.POST("/authenticate", s.user.Authenticate)
	}
}
