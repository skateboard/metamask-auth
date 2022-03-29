package session

import (
	"api/cache"
	"api/models"
	"api/routes/structs"
	"api/socket"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"os"
	"time"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func CreateSession(c *fiber.Ctx) error {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = false
	atClaims["expires"] = time.Now().Add(time.Hour * 1).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(structs.MessageStruct{
			Message: "Failed to generate JWT ticket.",
		})
	}

	ok := cache.Connection.AddSessionToCache(models.Session{
		SessionID: token,
	})
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(structs.MessageStruct{
			Message: "Failed to generate JWT ticket.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(structs.SuccessfulResponse{
		Data: token,
	})
}

func CheckSignature(c *fiber.Ctx) error {
	var newCheck NewCheckSignature
	if err := c.BodyParser(&newCheck); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(structs.MessageStruct{
			Message: "Failed to parse body.",
		})
	}
	sessionStorage := c.Locals("session").(*session.Session)
	sessionToken := sessionStorage.Get("session_token").(string)

	// Check if the token is valid

	socketConnection := socket.Socket.GetConnection(sessionToken)
	if socketConnection == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(structs.MessageStruct{
			Message: "Failed to find socket connection for this session!",
		})
	}

	newMessage := socket.Message{
		Type:    "AUTHENTICATED",
		Payload: User{
			Username: "CoolUser32",
		},
	}
	jsonBytes, err := json.Marshal(newMessage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(structs.MessageStruct{
			Message: "Failed to marshal message.",
		})
	}

	socketConnection.Send(jsonBytes)

	cache.Connection.DeleteSession(sessionToken)

	return c.Status(fiber.StatusOK).JSON(structs.SuccessfulResponse{
		Data: "OK",
	})
}