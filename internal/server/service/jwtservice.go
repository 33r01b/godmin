package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"godmin/config"
	"godmin/internal/dto"
	"godmin/internal/model"
	"godmin/internal/server/request"
	"godmin/internal/server/response"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
	"godmin/internal/throw"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

// JWTService is JWT authentication manager
type JWTService struct {
	store       *sqlstore.Store
	memoryStore *memorystore.Store
	config      *config.Jwt
}

// NewJwtService construct new JWTService
func NewJwtService(store *sqlstore.Store, memoryStore *memorystore.Store, jwtConfig *config.Jwt) *JWTService {
	return &JWTService{
		store:       store,
		memoryStore: memoryStore,
		config:      jwtConfig,
	}
}

// CreateToken build new JWT
func (s *JWTService) CreateToken(l *request.Login) (*response.Token, *throw.ResponseError) {
	u, err := s.store.User().FindByEmail(l.Email)
	if err != nil || !u.ComparePassword(l.Password) {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errIncorrectEmailOrPassword)
	}

	token, err := s.createToken(u.ID)
	if err != nil {
		return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
	}

	saveErr := s.memoryStore.Token().Create(u.ID, token)
	if saveErr != nil {
		return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
	}

	return &response.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// RefreshToken re-build JWT token
func (s *JWTService) RefreshToken(r *http.Request) (map[string]string, *throw.ResponseError) {
	req := &struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, throw.NewJWTError(http.StatusBadRequest, err)
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.RefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errors.New("refresh token expired"))
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, throw.NewJWTError(http.StatusUnauthorized, err)
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
		}

		//Delete the previous Access Token
		deleted, err := s.deleteToken(r)
		if err != nil || deleted == 0 {
			return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
		}

		//Delete the previous Refresh Token
		deleted, delErr := s.memoryStore.Token().Delete(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := s.createToken(userID)
		if createErr != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, createErr)
		}
		//save the tokens metadata to redis
		saveErr := s.memoryStore.Token().Create(userID, ts)
		if saveErr != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, saveErr)
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		return tokens, nil
	}

	return nil, throw.NewJWTError(http.StatusUnauthorized, errors.New("refresh expired"))
}

// Authenticate user by JWT token
func (s *JWTService) Authenticate(r *http.Request) (*model.User, *throw.ResponseError) {
	tokenAuth, errToken := s.extractTokenMetadata(r)
	if errToken != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	userID, errUserID := s.memoryStore.Token().Find(tokenAuth.AccessUUID)
	if errUserID != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}
	if userID != tokenAuth.UserID {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	u, errUser := s.store.User().Find(userID)
	if errUser != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	return u, nil
}

// Logout user
func (s *JWTService) Logout(r *http.Request) *throw.ResponseError {
	deleted, err := s.deleteToken(r)
	if err != nil || deleted == 0 {
		return throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	return nil
}

func (s *JWTService) createToken(userID uint64) (*dto.Token, error) {
	var err error
	token := &dto.Token{
		AccessUuid:          uuid.New().String(),
		RefreshUuid:         uuid.New().String(),
		RefreshTokenExpires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		AccessTokenExpires:  time.Now().Add(time.Minute * 15).Unix(),
	}

	// generate access token
	accessTokenClaims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": token.AccessUuid,
		"user_id":     userID,
		"exp":         token.AccessTokenExpires,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	token.AccessToken, err = accessToken.SignedString([]byte(s.config.AccessSecret))
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshTokenClaims := jwt.MapClaims{
		"refresh_uuid": token.RefreshUuid,
		"user_id":      userID,
		"exp":          token.RefreshTokenExpires,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(s.config.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *JWTService) deleteToken(r *http.Request) (int64, error) {
	tokenAuth, err := s.extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	var deleted int64
	deleted, err = s.memoryStore.Token().Delete(tokenAuth.AccessUUID)
	if err != nil || deleted == 0 {
		return 0, err
	}

	return deleted, nil
}

func (s *JWTService) extractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := s.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &accessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

func (s *JWTService) verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

type accessDetails struct {
	AccessUUID string
	UserID     uint64
}
