package handlers

import (
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/just-umyt/task/internal/models"
	"github.com/just-umyt/task/utils"
)

type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (models.User, error)
	UpdateUserToken(user *models.User) error
}

type AuthHandler struct {
	URepo UserRepo
}

func NewAuthHandler(userRepo UserRepo) *AuthHandler {
	return &AuthHandler{URepo: userRepo}
}

// UserSignUp method to create a new user.
func (authHandler *AuthHandler) UserSignUp(c *fiber.Ctx) error {
	//body parse
	signUp := models.SignUp{}

	if err := c.BodyParser(&signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//validate the response
	validate := validator.New()

	if err := validate.Struct(signUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	//Create a new user struct
	user := &models.User{}

	//hash the password
	userHashPwd, err := utils.GeneratePassword(signUp.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//set initialized default data for user:
	user.ID = uuid.New()
	user.Email = signUp.Email
	user.PasswordHash = userHashPwd
	user.CreatedAt = time.Now()

	//create new user in database
	if err := authHandler.URepo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//return StatusOk and new user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}

func (authHandler *AuthHandler) UserSignIn(c *fiber.Ctx) error {
	//get body to new struct
	signIn := &models.SignIn{}

	if err := c.BodyParser(&signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//validate
	validate := validator.New()

	if err := validate.Struct(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Get user by email
	foundedUser, err := authHandler.URepo.GetUserByEmail(signIn.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//compare hash
	compareUserPass := utils.ComparePassword(foundedUser.PasswordHash, signIn.Password)
	if !compareUserPass {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "wrong user email address or password",
		})
	}

	//Generate jwt tokens
	tokens, err := utils.CreateNewTokens(foundedUser.ID, c.IP())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Hash refresh token
	foundedUser.RefreshTokenHash = utils.NewHashedToken(tokens.Refresh)

	//Set hashed refresh token to db
	if err := authHandler.URepo.UpdateUserToken(&foundedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Encode Refresh token to base64
	tokens.Refresh = utils.EncodeToBase(tokens.Refresh)
	refreshTokensExpires, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	//Set refresh token to cookie
	c.Cookie(&fiber.Cookie{
		Name:    "refresh",
		Value:   tokens.Refresh,
		Expires: time.Now().Add(time.Hour * time.Duration(refreshTokensExpires)),
		// HTTPOnly: true,
		// Secure:   true,
		// SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"error":  false,
		"msg":    nil,
		"tokens": tokens,
	})
}
