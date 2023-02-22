package controller

import (
	"anara/entity"
	"anara/services"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	adminService services.AdminService
}

func NewAuthController(adminService services.AdminService) *AuthController {
	return &AuthController{
		adminService: adminService,
	}
}

// @Summary Login
// @Tags Admin Auth
// @Description Login with email and password
// @Param  input body entity.LoginReq true "login user input"
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.LoginResp
// @Failure 400 {object} entity.ErrRespController
// @Failure 404 {object} entity.ErrRespController
// @Router /admin/login [post]
func (a *AuthController) Login(c *fiber.Ctx) error {
	functionName := "Login"

	var input entity.LoginReq

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing login input, details %v", err),
		})
	}

	admin, err := a.adminService.GetAdminByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "invalid email or password",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["admin_id"] = admin.AdminID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on signing access token, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.LoginResp{
		AccessToken: t,
	})
}
