package main

import (
	"fmt"
	"test-be-kalbe/internal/infrastructure"
)

func main() {
	viperConfig := infrastructure.NewViper()
	log := infrastructure.NewLogger(viperConfig)
	db := infrastructure.NewDatabase(viperConfig, log)
	validate := infrastructure.NewValidator(viperConfig)
	app := infrastructure.NewFiber(viperConfig)

	infrastructure.Bootstrap(&infrastructure.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
