package middleware

import (
	"anara/entity"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

func accessTokenErrJWT(c *fiber.Ctx, err error) error {
	functionName := "accessTokenErrJWT"

	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
		SourceFunction: functionName,
		ErrMessage:     "invalid or expired jwt",
	})
}

// Guards a specific endpoint in the API.
func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler:  accessTokenErrJWT,
		SigningKey:    []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		SigningMethod: "HS256",
		TokenLookup:   "header:x-access-token",
	})
}

// Gets user data (their ID) from the JWT middleware. Should be executed after calling 'JWTMiddleware()'.
func GetDataFromJWT(c *fiber.Ctx) error {
	// Get userID from the previous route.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)
	adminId := claims["admin_id"].(float64)

	// Go to next.
	c.Locals("admin_id", adminId)
	return c.Next()
}
