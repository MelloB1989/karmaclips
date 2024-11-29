package auth

import (
	"fmt"
	"karmaclips/config"
	"karmaclips/database"
	"karmaclips/helpers/users"
	"karmaclips/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		fmt.Println("Error parsing request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	user, err := users.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		fmt.Println("Error getting user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email does not exist",
		})
	}
	if utils.CheckPasswordHash(req.Password, user.Password) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user.Email
		claims["uid"] = user.Id
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
		t, err := token.SignedString([]byte(config.NewConfig().JWTSecret))
		if err != nil {
			fmt.Println("Error signing token")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Error signing token",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": t,
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func Register(c *fiber.Ctx) error {
	userReq := new(RegisterRequest)
	if err := c.BodyParser(userReq); err != nil {
		fmt.Println("Error parsing request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		fmt.Println("Error hashing password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}
	userReq.Password = hashedPassword
	user := &database.Users{
		Email:    userReq.Email,
		Password: userReq.Password,
		Name:     userReq.Name,
		Id:       utils.GenerateID(),
	}
	_, err = users.CreateUser(user)
	if err != nil {
		fmt.Println("Error creating user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error creating user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
	})
}
