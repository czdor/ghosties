package services

import (
	"core/config"
	"core/internal/adapters/database/models"
	repositories "core/internal/ports"
	"core/types"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCharacterLimit = 72 // ? https://www.ory.sh/docs/troubleshooting/bcrypt-secret-length
)

var (
	ErrorInvalidCredentials        = types.ApiError("Invalid Email and/or Password")
	ErrorEmailNotRegistered        = types.ApiError("Email is not registered")
	ErrorFailedAccessTokenGen      = types.ApiError("Failed to generate access/refresh token")
	ErrorEmailAlreadyRegistered    = types.ApiError("Email already exists")
	ErrorUsernameAlreadyRegistered = types.ApiError("Username already exists")
	ErrorFailedToHashPassword      = types.ApiError("Failed to hash password")
	ErrorPasswordLenExceeded       = types.ApiError("Password must be no longer than 72 characters")
	ErrorSaveFailed                = types.ApiError("Failed to save")
)

type UserService struct {
	userRepo *repositories.UserRepoContext
}

func NewUserService(userRepo *repositories.UserRepoContext) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (ctx *UserService) Login(email string, password string) (int, types.ApiResponse) {
	user, err := ctx.userRepo.GetByEmail(email)
	if err != nil {
		return http.StatusBadRequest, ErrorEmailNotRegistered
	}

	if user.ID == 0 {
		return http.StatusBadRequest, ErrorInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return http.StatusBadRequest, ErrorEmailNotRegistered
	}

	// * Generate jwt access and refresh pair tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return http.StatusBadRequest, ErrorFailedAccessTokenGen
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return http.StatusBadRequest, ErrorFailedAccessTokenGen
	}

	// * Send a response
	return http.StatusOK, types.ApiResponse{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	}
}

func (ctx *UserService) Signup(email string, username string, password string) (int, types.ApiResponse) {
	if len(password) > BcryptCharacterLimit {
		return http.StatusBadRequest, ErrorPasswordLenExceeded
	}

	// * 2. Check if Email or Username is already stored
	exists := ctx.userRepo.ExistsEmail(email)
	if exists {
		return http.StatusBadRequest, ErrorEmailAlreadyRegistered
	}

	exists = ctx.userRepo.ExistsUsername(username)
	if exists {
		return http.StatusBadRequest, ErrorUsernameAlreadyRegistered
	}

	// * 3. Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusBadRequest, ErrorFailedToHashPassword
	}

	// * 4. Save user
	user := models.User{
		Email:    email,
		Username: strings.TrimSpace(username),
		Password: string(passwordHash),
	}

	if _, err := ctx.userRepo.Save(user); err != nil {
		return http.StatusBadRequest, ErrorSaveFailed
	}

	// * 5. Send a response
	return http.StatusCreated, types.ApiResponse{
		"success": true,
	}
}

func (ctx *UserService) RefreshToken(user models.User) (int, types.ApiResponse) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return http.StatusBadRequest, ErrorFailedAccessTokenGen
	}

	return http.StatusOK, types.ApiResponse{"accessToken": refreshTokenString}
}
