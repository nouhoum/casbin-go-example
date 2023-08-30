package database

import (
	"fmt"
	"log"

	"github.com/nouhoum/casbin-go-example/internal/model"
	"github.com/samber/do"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var models = []interface{}{
	&model.TodoItem{},
	&model.User{},
}

func New(i *do.Injector) (*gorm.DB, error) {
	log.Println("creating db object...")
	cfg := do.MustInvoke[*Config](i)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DSN(),
	}))

	if err != nil {
		return nil, err
	}

	log.Println("db object created...")
	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, fmt.Errorf("unable to perform auto migration: %v", err)
	}
	return db, nil
}
