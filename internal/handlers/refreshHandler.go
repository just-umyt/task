package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/just-umyt/task/internal/models"
	"github.com/just-umyt/task/utils"
)

type TokenRepository interface {
	GetUserById(id uuid.UUID) (models.User, error)
}

type TokenHandler struct {
	TokenRepo TokenRepository
}

func NewTokenHandler(tokenRepo TokenRepository) *TokenHandler {
	return &TokenHandler{TokenRepo: tokenRepo}
}

func (tokenrepo *TokenHandler) RefreshToken(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	basedRefreshToken := c.Cookies("refresh")

	//get claims from access token
	accessClaims, err := utils.ParseToken(accessToken[7:], os.Getenv("JWT_SECRET_KEY"))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Checking refresh token
	if basedRefreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Unauthorized",
		})
	}

	//decode from base
	refreshToken, err := utils.DecodeFromBase(basedRefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//get claims from refresh token
	refreshClaims, err := utils.ParseToken(refreshToken, os.Getenv("JWT_REFRESH_KEY"))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Checking access and refresh Id
	if accessClaims.ID != refreshClaims.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "токены Не связанные",
		})
	}

	//check user d in tokens
	if accessClaims.UserId != refreshClaims.UserId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Токены фальшивые. Авторизуйтесь заново",
		})
	}

	//find user by id
	foundedUser, err := tokenrepo.TokenRepo.GetUserById(refreshClaims.UserId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//compare with database refresh token
	if !utils.CompareToken(foundedUser.RefreshTokenHash, refreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "token compare is false",
		})
	}

	//compare ip address
	// ip = "another ip" //to check

	if c.IP() != refreshClaims.Ip && c.IP() != accessClaims.Ip {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "someone tries login to your account",
		})
	}

	//give a new acces token
	newAccessToken, err := utils.NewAccessToken(*refreshClaims)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"tokens": fiber.Map{
			"refresh Token":    refreshToken,
			"new Access Token": newAccessToken,
		},
	})
}
