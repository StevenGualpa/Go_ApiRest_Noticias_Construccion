package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"noticias_uteq/repositories"
	"noticias_uteq/routers"
)

func main() {
	app := fiber.New()
	viper.SetDefault("port", "3000")
	viper.AutomaticEnv()

	repositories.InitDB()
	repositories.Migrate()

	// Configuramos las rutas
	routers.SetupRouter(app)

	err := app.Listen(":" + viper.GetString("port"))
	if err != nil {
		log.Fatalln(err)
	}
}
