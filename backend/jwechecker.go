package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-jose/go-jose/v3"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/hkdf"
)

// should me stored in env variable
const NEXTAUTH_SECRET = "LnydL9vQ0oAOh8otMc5dDaOtSmHIPNKyKNlB/y7br5M="

func getDerivedEncryptionKey() []byte {
	hkdf := hkdf.New(sha256.New, []byte(NEXTAUTH_SECRET), nil, []byte("NextAuth.js Generated Encryption Key"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		panic(err)
	}
	return key
}

var secret = getDerivedEncryptionKey()

type JweConfig struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx, cfg *JweConfig) (*jwt.MapClaims, error)
	Secret       []byte
	Expiry       int64
	ContextKey   string
}

var ConfigDefault = JweConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
	Secret:       secret,
	Expiry:       60, //60 * 60 * 24 * 30,
	ContextKey:   "jwtClaims",
}

func Decode(c *fiber.Ctx, cfg *JweConfig) (*jwt.MapClaims, error) {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return nil, errors.New("authorization header is required")
	}

	tokenString := strings.Split(authHeader, "Bearer ")

	if len(tokenString) < 2 {
		return nil, errors.New("authroization Bearer is required")
	}

	jwtEncrypted, err := jose.ParseEncrypted(tokenString[1])
	if err != nil {
		return nil, errors.New("error parsing encrypted token")
	}

	decryptedJWE, err := jwtEncrypted.Decrypt(cfg.Secret)
	if err != nil {
		return nil, errors.New("error decrypting token")
	}
	token := string(decryptedJWE)

	fmt.Println(token)

	jwtJson := parseMap(token)

	// fmt.Println(jwtJson)

	claims := jwt.MapClaims{}

	for k, value := range jwtJson {
		claims[k] = value
	}

	// fmt.Println(claims)

	err = claims.Valid()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	return &claims, nil
}

func parseMap(tokenString string) map[string]interface{} {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(tokenString), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func Unauthorized(c *fiber.Ctx) error {
	fmt.Println("Unauthorized")
	return c.SendStatus(fiber.StatusUnauthorized)
}

func jweConfigDefault(config ...JweConfig) JweConfig {
	if len(config) < 1 {
		return ConfigDefault
	}
	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}
	// Set default secret if not passed
	if len(cfg.Secret) < 1 {
		cfg.Secret = ConfigDefault.Secret
	}

	// Set default expiry if not passed
	if cfg.Expiry == 0 {
		cfg.Expiry = ConfigDefault.Expiry
	}

	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}

	if cfg.Decode == nil {
		cfg.Decode = Decode
	}
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = Unauthorized
	}
	return cfg
}

func NewJweConfig(config JweConfig) fiber.Handler {
	cfg := jweConfigDefault(config)

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			fmt.Println("Middleware was skipped")
			return c.Next()
		}
		fmt.Println("Middleware was run")

		claims, err := cfg.Decode(c, &cfg)

		if err == nil {
			c.Locals(cfg.ContextKey, *claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}
