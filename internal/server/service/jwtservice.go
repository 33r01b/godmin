package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"godmin/config"
	"godmin/internal/dto"
	"godmin/internal/model"
	"godmin/internal/server/request"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
	"godmin/internal/throw"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type JWTService struct {
	store       *sqlstore.Store
	memoryStore *memorystore.Store
	config      *config.Jwt
}

func NewJwtService(store *sqlstore.Store, memoryStore *memorystore.Store, jwtConfig *config.Jwt) *JWTService {
	return &JWTService{
		store:       store,
		memoryStore: memoryStore,
		config:      jwtConfig,
	}
}

func (s *JWTService) CreateToken(l *request.Login) (map[string]string, *throw.ResponseError) {
	u, err := s.store.User().FindByEmail(l.Email)
	if err != nil || !u.ComparePassword(l.Password) {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errIncorrectEmailOrPassword)
	}

	var ts *dto.Token
	ts, err = s.createToken(u.ID)
	if err != nil {
		return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
	}

	saveErr := s.memoryStore.Token().Create(u.ID, ts)
	if saveErr != nil {
		return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	return tokens, nil
}

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
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, err)
		}

		//Delete the previous Access Token
		deleted, err := s.deleteToken(r)
		if err != nil || deleted == 0 {
			return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
		}

		//Delete the previous Refresh Token
		deleted, delErr := s.memoryStore.Token().Delete(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := s.createToken(userId)
		if createErr != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, createErr)
		}
		//save the tokens metadata to redis
		saveErr := s.memoryStore.Token().Create(userId, ts)
		if saveErr != nil {
			return nil, throw.NewJWTError(http.StatusUnprocessableEntity, saveErr)
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		return tokens, nil
	} else {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errors.New("refresh expired"))
	}
}

func (s *JWTService) Authenticate(r *http.Request) (*model.User, *throw.ResponseError) {
	tokenAuth, errToken := s.extractTokenMetadata(r)
	if errToken != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	userId, errUserId := s.memoryStore.Token().Find(tokenAuth.AccessUuid)
	if errUserId != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}
	if userId != tokenAuth.UserId {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	u, errUser := s.store.User().Find(userId)
	if errUser != nil {
		return nil, throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	return u, nil
}

func (s *JWTService) Logout(r *http.Request) *throw.ResponseError {
	deleted, err := s.deleteToken(r)
	if err != nil || deleted == 0 {
		return throw.NewJWTError(http.StatusUnauthorized, errNotAuthenticated)
	}

	return nil
}

func (s *JWTService) createToken(userId uint64) (*dto.Token, error) {
	t := &dto.Token{}
	t.AccessUuid = uuid.New().String()
	t.RefreshUuid = uuid.New().String()
	t.RefreshTokenExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.AccessTokenExpires = time.Now().Add(time.Minute * 15).Unix()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = t.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = t.AccessTokenExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(s.config.AccessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = t.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = t.RefreshTokenExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(s.config.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *JWTService) deleteToken(r *http.Request) (int64, error) {
	tokenAuth, err := s.extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	var deleted int64
	deleted, err = s.memoryStore.Token().Delete(tokenAuth.AccessUuid)
	if err != nil || deleted == 0 {
		return 0, err
	}

	return deleted, nil
}

func (s *JWTService) extractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := s.verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
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

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}
