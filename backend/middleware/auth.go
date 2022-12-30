package middleware

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/go-jose/go-jose/v3"
	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/hkdf"
)

var nextauth_secret string

type TokenClaim struct {
	Name              string
	Email             string
	Picture           string
	Sub               string
	AccessToken       string
	Id                int64
	ProviderAccountId string
	Provider          string
	Iat               int64
	Exp               int64
	Jti               string
	Role              string
}

type AuthConfig struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx, cfg *AuthConfig) (*jwt.MapClaims, error)
	Secret       []byte
	Expiry       int64
	ContextKey   string
}

var defaultConfig = AuthConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
	Secret:       nil,
	Expiry:       60 * 60, //60 * 60 * 24 * 30,
	ContextKey:   "user",
}

func Auth(config AuthConfig) fiber.Handler {
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

func decode(c *fiber.Ctx, cfg *AuthConfig) (*jwt.MapClaims, error) {
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
	claims := parseClaim(string(decryptedJWE))

	/* check claims jwt is valid */
	err = claims.Valid()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	return &claims, nil
}

func parseClaim(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	token := TokenClaim{}
	err := json.Unmarshal([]byte(tokenString), &token)
	if err != nil {
		panic(err)
	}
	values := reflect.ValueOf(token)
	typesOf := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).Interface() != nil {
			claims[typesOf.Field(i).Name] = values.Field(i).Interface()
		}
	}
	return claims
}

func getDerivedEncryptionKey() []byte {
	if nextauth_secret == "" {
		nextauth_secret = envVariable("NEXTAUTH_SECRET")
	}
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

func unauthorized(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusUnauthorized)
}

func configDefaults(config ...AuthConfig) AuthConfig {
	if len(config) < 1 {
		return defaultConfig
	}
	// Override default config
	cfg := config[0]

	// initialize default secret
	if defaultConfig.Secret == nil {
		defaultConfig.Secret = getDerivedEncryptionKey()
	}

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
