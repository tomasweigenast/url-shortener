package services

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/mapper"
	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/utils"
)

type authService struct {
	jwtPass []byte
	jwtTtl  time.Duration
}

type TokenClaims struct {
	Sub   uint32 `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

var authServiceSingleton *authService

func AuthService() *authService {
	if authServiceSingleton == nil {
		jwtPass := []byte(os.Getenv("JWT_PASS"))
		if len(jwtPass) == 0 {
			log.Fatalf("JWT_PASS environment variable is missing")
		}

		jwtTtl, err := strconv.ParseUint(os.Getenv("JWT_TTL"), 10, 32)
		if err != nil {
			log.Fatalf("unable to parse JWT_TTL enviroment variable: %s", err)
		}

		authServiceSingleton = &authService{
			jwtPass: jwtPass,
			jwtTtl:  time.Duration(jwtTtl) * time.Second,
		}
	}

	return authServiceSingleton
}

// SignInUser finds an user, verifies its password and if suceeded, generates a new token and returns it
func (as *authService) SignInUser(ctx context.Context, email, password string) (*models.User, *models.AccessToken, error) {
	user, err := database.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if user.IsDeleted {
		return nil, nil, utils.ErrNotFound
	}

	if !utils.VerifyPassword(user.PasswordHash, password) {
		return nil, nil, errors.New("wrong-credentials")
	}

	if user.IsDisabled {
		return nil, nil, errors.New("user-disabled")
	}

	token, err := as.SignToken(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return mapper.MapUser(user), token, nil
}

// SignToken creates a new session and an access token. It does not verifies user password nor if its enabled.
func (as *authService) SignToken(ctx context.Context, user *entities.User) (*models.AccessToken, error) {
	now := time.Now()
	expiration := now.Add(as.jwtTtl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "url-shortener.example",
			Subject:   "url-shortener",
		},
		Sub:   user.Id,
		Email: user.Email,
		Name:  user.Name,
	})

	tokenString, err := token.SignedString(as.jwtPass)
	if err != nil {
		return nil, err
	}

	refreshToken := utils.RandomString()
	refreshTokenExpiration := now.Add(as.jwtTtl * 2)

	err = database.InsertSession(ctx, &entities.UserSession{
		Id:           utils.RandomId(),
		Uid:          user.Id,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenExpiration,
	})

	if err != nil {
		return nil, err
	}

	return &models.AccessToken{
		Token:        tokenString,
		Expires:      expiration,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken validates if the given bearer token is valid and not expired
func (as *authService) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	claims := TokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return as.jwtPass, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token-invalid")
	}

	return &claims, nil
}
