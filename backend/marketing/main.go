package main

import (
	"context"
	"log"
	"marketing/handlers"
	"marketing/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main() {
    ctx := context.Background()
    services.Init()
    defer services.MongoClient.Disconnect(ctx)
    defer services.RabbitMQChannel.Close()
    defer services.RabbitMQConn.Close()

    app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
	}))
    
    app.Get("/minhas-inscricoes", handlers.MySubscriptions)
    app.Post("/inscrever", handlers.Subscribe)
    app.Post("/promocao", handlers.CreatePromotion)
    app.Post("/cancelar", handlers.CancelSubscription)
    
    log.Fatal(app.Listen(":3004"))
}