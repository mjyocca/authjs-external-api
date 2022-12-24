package main

import (
	"fmt"
	"io"

	"crypto/sha256"

	fiber "github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/hkdf"
)

const NEXTAUTH_SECRET = "LnydL9vQ0oAOh8otMc5dDaOtSmHIPNKyKNlB/y7br5M="

func getDerivedEncryptionKey() []byte {
	hkdf := hkdf.New(sha256.New, []byte(NEXTAUTH_SECRET), nil, []byte("NextAuth.js Generated Encryption Key"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		panic(err)
	}
	return key
}

func main() {
	app := fiber.New()

	encryptionSecret := getDerivedEncryptionKey()
	app.Use(jwtware.New(jwtware.Config{
		// SigningKey:    []byte(NEXTAUTH_SECRET),
		SigningKey:    encryptionSecret,
		TokenLookup:   "header:authorization",
		AuthScheme:    "bearer",
		SigningMethod: "HS256",
		ContextKey:    "user",
		// Claims:        jwt.MapClaims{},
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Locals("user"))
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		accessToken := claims["accessToken"].(string)
		// return c.SendString("Hello, World 👋!")
		return c.JSON(fiber.Map{
			"name":        name,
			"accessToken": accessToken,
		})
	})

	fmt.Printf("NEXTAUTH_SECRET %+v", NEXTAUTH_SECRET)

	app.Listen(":8000")
}
