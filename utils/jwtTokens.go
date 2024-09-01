package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
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
	UserId string `json:"user_id"`
	Ip     string `json:"ip"`
	jwt.RegisteredClaims
}

// Create a new pair of tokens
func CreateNewTokens(id uuid.UUID, ip string) (*Tokens, error) {
	jti := strconv.Itoa(rand.Intn(1000000000))

	//new access token
	accessToken, err := NewAccessToken(id, ip, jti)
	if err != nil {
		return nil, err
	}

	//new refresh token
	refreshToken, err := NewRefreshToken(id, ip, jti)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func NewAccessToken(id uuid.UUID, ip string, jti string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	type MyCustomClaims struct {
		UserId string `json:"user_id"`
		Ip     string `json:"ip"`
		jwt.RegisteredClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		id.String(),
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutesCount) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func NewRefreshToken(id uuid.UUID, ip string, jti string) (string, error) {
	secret := []byte(os.Getenv("JWT_REFRESH_KEY"))
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	// Create the Claims
	claims := MyCustomClaims{
		id.String(),
		ip,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hoursCount) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseRefreshToken(refresh string) (uuid.UUID, string, string, error) {

	token, err := jwt.ParseWithClaims(refresh, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_KEY")), nil
	})

	if err != nil {
		return uuid.Nil, "", "", err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		userId, _ := uuid.Parse(claims.UserId)
		return userId, claims.Ip, claims.ID, nil
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}

	return uuid.Nil, "", "", nil
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
	refreshT, err := base64.StdEncoding.DecodeString(based[7:])
	if err != nil {
		return "", nil
	}

	return string(refreshT), nil
}
