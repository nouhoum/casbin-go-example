package cmd

import (
	"fmt"
	"log"

	"github.com/nouhoum/casbin-go-example/internal/database"
	"github.com/nouhoum/casbin-go-example/internal/handler"
	"github.com/nouhoum/casbin-go-example/internal/server"
	"github.com/nouhoum/casbin-go-example/internal/service"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "casbin-go-example",
	Short: "Casbin Go Example",
	Run:   runServer,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("%v", err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	injector := do.New()
	doInjection(injector)

	if err := do.MustInvoke[service.Role](injector).InitRoles(); err != nil {
		log.Fatalf("failed to initialize roles. err=%v", err)
	}

	if err := do.MustInvoke[service.User](injector).InitUsers(); err != nil {
		log.Fatalf("failed to initialize users. err=%v", err)
	}

	server := do.MustInvoke[*server.Server](injector)
	if err := server.Run(); err != nil {
		log.Fatal("error server closed", err)
	}

	log.Println("server exiting...")
}

func doInjection(injector *do.Injector) {
	do.Provide(injector, database.NewConfig)
	do.Provide(injector, database.New)
	do.Provide(injector, handler.NewTodo)
	do.Provide(injector, handler.NewUser)
	do.Provide(injector, server.New)
	do.Provide(injector, server.NewConfig)
	do.Provide(injector, server.NewEngine)
	do.Provide(injector, service.NewAuthMiddleware)
	do.Provide(injector, service.NewCasbinConfig)
	do.Provide(injector, service.NewCasbinEnforcer)
	do.Provide(injector, service.NewJWTConfig)
	do.Provide(injector, service.NewPolicy)
	do.Provide(injector, service.NewRole)
	do.Provide(injector, service.NewTodo)
	do.Provide(injector, service.NewUser)
}

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8081)

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
