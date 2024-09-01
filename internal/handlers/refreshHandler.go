package handlers

import (
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
	basedRefreshToken := c.Get("Authorization")

	//decode from base
	refreshToken, err := utils.DecodeFromBase(basedRefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//get user id, ip, jit from refresh token
	userId, ip, jti, err := utils.ParseRefreshToken(refreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//find user by id
	foundedUser, err := tokenrepo.TokenRepo.GetUserById(userId)
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

	if c.IP() != ip {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "someone tries login to your account",
		})
	}

	//give a new acces token
	newAccessToken, err := utils.NewAccessToken(userId, c.IP(), jti)
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
