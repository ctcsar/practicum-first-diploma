package auth

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/ctcsar/practicum-first-diploma/internal/database"
	"github.com/ctcsar/practicum-first-diploma/internal/logger"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"password"`
	RegisterDate string    `json:"registerDate,omitempty"`
}

var JWT_KEY string

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}
	JWT_KEY = os.Getenv("JWT_SECRET")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(login string, password string) (*User, error) {
	id := uuid.New()
	newPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           id,
		Login:        login,
		Password:     newPassword,
		RegisterDate: time.Now().String(),
	}, nil

}

func AuthUser(login string, password string, db *sql.DB) (string, error) {
	var uuid uuid.UUID
	var pass string
	fmt.Printf("JWT_KEY: %s\n", JWT_KEY)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}
	uuid, pass, err = database.CheckUser(login, hashedPassword, db)
	if err != nil {
		return "", err
	}
	if CheckPasswordHash(password, pass) {
		token, err := GenerateToken(uuid, login)
		if err != nil {
			return "", err
		}
		return token, nil
	} else {
		return "", fmt.Errorf("wrong login or password")
	}

}

func GenerateToken(id uuid.UUID, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(JWT_KEY))
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
