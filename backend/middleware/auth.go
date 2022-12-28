package middleware

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-jose/go-jose/v3"
	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/hkdf"
)

// should me stored in env variable
var nextauth_secret = envVariable("NEXTAUTH_SECRET")
var secret = getDerivedEncryptionKey()

func getDerivedEncryptionKey() []byte {
	hkdf := hkdf.New(sha256.New, []byte(nextauth_secret), nil, []byte("NextAuth.js Generated Encryption Key"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		panic(err)
	}
	return key
}

func envVariable(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Error getting env variable %+s", key)
	}
	return val
}

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx, cfg *Config) (*jwt.MapClaims, error)
	Secret       []byte
	Expiry       int64
	ContextKey   string
}

var defaultConfig = Config{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
	Secret:       secret,
	Expiry:       60, //60 * 60 * 24 * 30,
	ContextKey:   "jwtClaims",
}

func decode(c *fiber.Ctx, cfg *Config) (*jwt.MapClaims, error) {
	authHeader := c.Get("Authorization")

	/* check request is valid */
	if authHeader == "" {
		return nil, errors.New("authorization header is required")
	}

	tokenString := strings.Split(authHeader, "Bearer ")

	if len(tokenString) < 2 {
		return nil, errors.New("authroization Bearer is required")
	}
	/* check request is valid */

	/* parse & decrypt jwe */
	jwtEncrypted, err := jose.ParseEncrypted(tokenString[1])
	if err != nil {
		return nil, errors.New("error parsing encrypted token")
	}

	decryptedJWE, err := jwtEncrypted.Decrypt(cfg.Secret)
	if err != nil {
		return nil, errors.New("error decrypting token")
	}
	/* parse & decrypt jwe */

	/* convert decrypted jwt to claims */
	token := string(decryptedJWE)
	jwtJson := parseMap(token)
	claims := jwt.MapClaims{}
	for k, value := range jwtJson {
		claims[k] = value
	}

	/* check claims jwt is valid */
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

func unauthorized(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusUnauthorized)
}

func configDefaults(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}
	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Filter == nil {
		cfg.Filter = defaultConfig.Filter
	}
	// Set default secret if not passed
	if len(cfg.Secret) < 1 {
		cfg.Secret = defaultConfig.Secret
	}

	// Set default expiry if not passed
	if cfg.Expiry == 0 {
		cfg.Expiry = defaultConfig.Expiry
	}

	// Set default context key to retrieve claims via, c.Locals(ContextKey)
	if cfg.ContextKey == "" {
		cfg.ContextKey = defaultConfig.ContextKey
	}

	// Set decode function
	if cfg.Decode == nil {
		cfg.Decode = decode
	}
	// Set unauthorized handler when token is invalid/expired
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = unauthorized
	}
	return cfg
}

func Auth(config Config) fiber.Handler {
	cfg := configDefaults(config)

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
