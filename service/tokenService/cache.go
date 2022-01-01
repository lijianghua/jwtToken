package tokenService

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"jwtToken/config"
	"net/http"
	"time"
)

type TokenStorage interface {
	Get(key string) (string, error)
	Set(key, value string, expiresIn time.Duration) error
	Del(key string) (int64, error)
}

type Manager struct {
	storage TokenStorage
	opts    *config.JwtConfig
}

func NewService(storage TokenStorage, opts *config.JwtConfig) *Manager {
	return &Manager{
		storage: storage,
		opts:    opts,
	}
}

//NewToken: 创建Token
func (m Manager) NewToken(userID string) (*Token, error) {

	accessDuration, err := time.ParseDuration(m.opts.AccessTokenDuration)
	if err != nil {
		accessDuration = 15 * time.Minute
	}
	refreshDuration, err := time.ParseDuration(m.opts.RefreshTokenDuration)
	if err != nil {
		refreshDuration = 7 * 24 * time.Hour
	}
	now := time.Now()

	token := &Token{
		AccessID:         uuid.NewV4().String(),
		AccessCreateAt:   now,
		AccessExpiresIn:  accessDuration,
		RefreshCreateAt:  now,
		RefreshExpiresIn: refreshDuration,
	}
	token.RefreshID = token.AccessID + "++" + userID

	// Create access token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["access_id"] = token.AccessID
	claims["exp"] = token.AccessCreateAt.Add(token.AccessExpiresIn).Unix()
	//&Claims{
	//	StandardClaims: jwt.StandardClaims{
	//		Subject: username,
	//		// In JWT, the expiry time is expressed as unix milliseconds
	//		ExpiresAt: token.AccessCreateAt.Add(token.AccessExpiresIn).Unix(),
	//	},
	//}
	token.Access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.opts.JwtAccessSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["refresh_id"] = token.RefreshID
	rtClaims["exp"] = token.RefreshCreateAt.Add(token.RefreshExpiresIn).Unix()
	token.Refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString([]byte(m.opts.JwtRefreshSecret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m Manager) ParseRefreshToken(refreshToken string) (*jwt.Token, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.opts.JwtRefreshSecret), nil
	})

	return token, err
}

func (m Manager) CacheAuth(userID string, token *Token) error {

	err := m.storage.Set(token.AccessID, userID, token.AccessExpiresIn)
	if err != nil {
		return err
	}
	err = m.storage.Set(token.RefreshID, userID, token.RefreshExpiresIn)
	return err
}

func (m Manager) FetchAuth(authD *AccessDetails) (string, error) {
	userid, err := m.storage.Get(authD.AccessID)
	if err != nil {
		return "", err
	}
	return userid, nil
}

func (m Manager) DeleteCacheAuth(givenUuid string) (int64, error) {
	deleted, err := m.storage.Del(givenUuid)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (m Manager) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.AccessID, authD.UserID)
	//delete access token
	deletedAt, err := m.storage.Del(authD.AccessID)
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := m.storage.Del(refreshUuid)
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("redis delete token failed")
	}
	return nil
}

func (m Manager) VerifyToken(r *http.Request) (*jwt.Token, error) {

	access, ok := ExtractToken(r)
	if !ok {
		return nil, errors.New("Can't find access token from Bearer Authorization header")
	}

	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.opts.JwtAccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
func (m Manager) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := m.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessID, ok := claims["access_id"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessID: accessID,
			UserID:   userID,
		}, nil
	}
	return nil, err
}
