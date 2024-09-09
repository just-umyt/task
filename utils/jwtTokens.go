package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	Access  string
	Refresh string
}

type MyCustomClaims struct {
	UserId uuid.UUID `json:"user_id"`
	Ip     string    `json:"ip"`
	jwt.RegisteredClaims
}

// Create a new pair of tokens
func CreateNewTokens(id uuid.UUID, ip string) (*Tokens, error) {
	jti := strconv.Itoa(rand.Intn(1000000000))

	//new access token
	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	accessClaims := MyCustomClaims{
		id,
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesCount) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	accessToken, err := NewAccessToken(accessClaims)
	if err != nil {
		return nil, err
	}

	//new refresh token
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	refreshClaims := MyCustomClaims{
		id,
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hoursCount) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	refreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func NewAccessToken(claims MyCustomClaims) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func NewRefreshToken(claims MyCustomClaims) (string, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseToken(tokenString, secret string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Convert token to hash
func NewHashedToken(token string) string {
	hash := sha256.New()

	hash.Write([]byte(token))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Compare hashed token
func CompareToken(hashedToken, inputToken string) bool {
	hash := sha256.New()

	hash.Write([]byte(inputToken))

	inputStr := fmt.Sprintf("%x", hash.Sum(nil))

	return hashedToken == inputStr
}

// Standart encoding to base64
func EncodeToBase(token string) string {
	data := []byte(token)
	return base64.StdEncoding.EncodeToString(data)
}

// Decode based token to string
func DecodeFromBase(based string) (string, error) {
	refreshT, err := base64.StdEncoding.DecodeString(based)
	if err != nil {
		return "", nil
	}

	return string(refreshT), nil
}
