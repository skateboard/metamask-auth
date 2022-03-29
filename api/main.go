package main

import (
	"api/middleware"
	"api/routes/session"
	"api/socket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main()  {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	port := os.Getenv("PORT")

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	apiGroup := app.Group("/v1")

	sessionGroup := apiGroup.Group("/session")
	sessionGroup.Post("/", session.CreateSession)
	sessionGroup.Post("/check_signature", middleware.SessionMiddleware, session.CheckSignature)

	socket.Initialize()
	wsGroup := app.Group("/v1/ws")
	wsGroup.Use("/", middleware.WSUpgrade)
	wsGroup.Get("/:sessionToken", websocket.New(socket.Socket.Handler))

	log.Println("API Listening on :" + port)
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start app! %s", err)
		return
	}
}
