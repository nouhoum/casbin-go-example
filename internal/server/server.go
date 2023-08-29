package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	}

	s.addRoutes()
	return &s, nil
}

func (s *Server) Run() error {
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
	//TODO: Add routes here.
}
