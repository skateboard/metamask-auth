package middleware

import (
	"api/cache"
	"api/routes/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SessionMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	store := session.New()
	s, _ := store.Get(c)

	if authHeader != "" {
		stored := cache.Connection.GetSession(authHeader)
		if stored == nil {
			return c.Status(401).JSON(structs.MessageStruct{
				Message: "Not a valid session!",
			})
		}

		s.Set("session_token", authHeader)
	} else {
		return c.Status(401).JSON(structs.MessageStruct{
			Message: "Failed to find 'Authorization' header",
		})
	}

	defer s.Save()
	c.Locals("session", s)
	return c.Next()
}