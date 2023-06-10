package service

import (
	"errors"
	"fmt"
	InternalRepository "go-api/repository/internal-repo"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTService is a contract of what jwtService can do
type JWTService interface {
	GenerateToken(userID string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type jwtCustomClaim2 struct {
	jwt.StandardClaims
}

type jwtService struct {
	secretKey      string
	issuer         string
	userRepository InternalRepository.UserRepository
}

// NewJWTService method is creates a new instance of JWTService
func NewJWTService(userRepo InternalRepository.UserRepository) JWTService {
	return &jwtService{
		issuer:         "x",
		secretKey:      getSecretKey(),
		userRepository: userRepo,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "crmBackend"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(UserID string) string {
	claims := &jwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(1 * time.Hour)).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {

	if tokenMatch := j.userRepository.FindToken(token); tokenMatch == "" {
		return nil, errors.New("Token not match")
	}

	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
