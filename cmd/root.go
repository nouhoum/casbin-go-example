package cmd

import (
	"log"

	"github.com/nouhoum/casbin-go-example/internal/server"
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

	server := do.MustInvoke[*server.Server](injector)

	if err := server.Run(); err != nil {
		log.Fatal("error server closed", err)
	}

	log.Println("server exiting...")
}

func doInjection(injector *do.Injector) {
	do.Provide(injector, server.New)
	do.Provide(injector, server.NewConfig)
	do.Provide(injector, server.NewEngine)
}

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8081)
}