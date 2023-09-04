package service

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nouhoum/casbin-go-example/api"
	"github.com/samber/do"
	"github.com/spf13/viper"
)

const (
	UserIDKey = "X-CURRENT-USER-ID"
)

type JWTConfig struct {
	realm string
	key   string
}

// NewJWTConfig constructs a new auth config
func NewJWTConfig(i *do.Injector) (*JWTConfig, error) {
	return &JWTConfig{
		realm: viper.GetString("JWT_REALM"),
		key:   viper.GetString("JWT_KEY"),
	}, nil
}

func NewAuthMiddleware(i *do.Injector) (*jwt.GinJWTMiddleware, error) {
	cfg := do.MustInvoke[*JWTConfig](i)
	svc := do.MustInvoke[User](i)
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            cfg.realm,
		Key:              []byte(cfg.realm),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		SigningAlgorithm: "HS256",
		IdentityKey:      UserIDKey,
		IdentityHandler:  idHandler,
		Authenticator:    authenticator(svc),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(api.AuthenticatedUser); ok {
				return jwt.MapClaims{
					"email":    v.Email,
					"id":       v.ID,
					"identity": v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*api.AuthenticatedUser); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return nil, err
	}

	return authMiddleware, nil
}

func authenticator(svc User) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var authReq api.AuthRequest
		if err := c.ShouldBind(&authReq); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}
		user, err := svc.Authenticate(c.Request.Context(), authReq.Email, authReq.Password)
		if err != nil {
			if err == ErrNoSuchUser {
				return nil, jwt.ErrFailedAuthentication
			}
			return nil, err
		}

		return api.AuthenticatedUser{ID: user.ID, Email: user.Email}, nil
	}
}

func idHandler(c *gin.Context) interface{} {
	return Authenticated(c)
}

// Authenticated extracts authenticated user from the request
func Authenticated(c *gin.Context) *api.AuthenticatedUser {
	claims := jwt.ExtractClaims(c)

	if claims != nil && claims["id"] != nil && claims["email"] != nil {
		return &api.AuthenticatedUser{
			ID:    int(claims["id"].(float64)),
			Email: claims["email"].(string),
		}
	}
	return nil
}
