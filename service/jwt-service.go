package service

import (
	"fmt"
	"go-api/dto/request/authRequestDTO"
	"go-api/entity"
	"go-api/repository"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// JWTService is a contract of what jwtService can do
type JWTService interface {
	GenerateToken(userID string) string
	GenerateToken2(subject string, request authRequestDTO.AuthRequest) string
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
	secretKey         string
	issuer            string
	jwtHistRepository repository.JwtHistRepository
}

// NewJWTService method is creates a new instance of JWTService
func NewJWTService(jwtHistRepo repository.JwtHistRepository) JWTService {
	return &jwtService{
		issuer:            "taufiq",
		secretKey:         getSecretKey(),
		jwtHistRepository: jwtHistRepo,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "taufiq"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(UserID string) string {
	claims := &jwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
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

func (j *jwtService) GenerateToken2(subject string, request authRequestDTO.AuthRequest) string {
	claims := &jwtCustomClaim2{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}

	uid := uuid.New()

	var data entity.JwtHistGo
	data.Email = request.Email
	data.Id = uid.String()
	data.Jwt = t

	j.jwtHistRepository.Save(data)

	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
