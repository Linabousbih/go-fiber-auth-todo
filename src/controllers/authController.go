package controllers

import (
	"fiberTodo/models"
	"fiberTodo/src/database"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getRequest(c fiber.Ctx) (Request, error) {

	var body Request

	if err := c.Bind().Body(&body); err != nil {
		return Request{}, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	return body, nil
}

func RegisterUser(c fiber.Ctx) error {

	body, err := getRequest(c)
	if err != nil {
		return err
	}

	var existingUser models.User

	err = database.DB.Collection("users").FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(bson.M{"error": "User already exists"})

	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	userDoc := bson.M{
		"email":    body.Email,
		"password": string(hashedPassword),
	}

	_, err = database.DB.Collection("users").InsertOne(c.Context(), userDoc)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func loginUser(c fiber.Ctx) error {
	body, err := getRequest(c)
	if err != nil {
		return err
	}
	var user models.User

	err = database.DB.Collection("users").FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Passowrd"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	t, err := token.SignedString(jwtSecret)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"message": "Logged in successfully"})
}
