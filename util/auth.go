package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"jwtToken/cache/redis"
	"jwtToken/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Bcrypt hash password
func EncodePass(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//Bcrypt verify password
func VerifyPass(encodedPwd string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPwd), []byte(pwd))
	return err == nil
}

type Token struct {
	AccessID        string        `json:"-"`
	Access          string        `json:"access-token"`
	AccessCreateAt  time.Time     `json:"access-token-create-at"`
	AccessExpiresIn time.Duration `json:"access-token-expires-in"`

	RefreshID        string        `json:"-"`
	Refresh          string        `json:"refresh-token"`
	RefreshCreateAt  time.Time     `json:"refresh-token-create-at"`
	RefreshExpiresIn time.Duration `json:"refresh-token-expires-in"`
}

type AccessDetails struct {
	AccessID string
	UserID   int64
}

//NewToken: 创建Token
func NewToken(userID int64) (*Token, error) {
	cfg := config.Cfg.Jwt
	accessDuration, err := time.ParseDuration(cfg.AccessTokenDuration)
	if err != nil {
		accessDuration = 15 * time.Minute
	}
	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenDuration)
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
	token.RefreshID = token.AccessID + "++" + strconv.Itoa(int(userID))

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
	token.Access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JwtAccessSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["refresh_id"] = token.RefreshID
	rtClaims["exp"] = token.RefreshCreateAt.Add(token.RefreshExpiresIn).Unix()
	token.Refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString([]byte(cfg.JwtRefreshSecret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

//ExtractToken: Extract access token from Bearer Authorization header
func ExtractToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	return token, token != ""
}

func (m *Token) CacheAuth(userID int64) error {
	var err error

	err = redis.RedisClient().Set(m.AccessID, strconv.Itoa(int(userID)), m.AccessExpiresIn).Err()
	if err != nil {
		return err
	}
	err = redis.RedisClient().Set(m.RefreshID, strconv.Itoa(int(userID)), m.RefreshExpiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

func FetchAuth(authD *AccessDetails) (int64, error) {
	userid, err := redis.RedisClient().Get(authD.AccessID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func DeleteCacheAuth(givenUuid string) (int64, error) {
	deleted, err := redis.RedisClient().Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessID, authD.UserID)
	//delete access token
	deletedAt, err := redis.RedisClient().Del(authD.AccessID).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := redis.RedisClient().Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("redis delete token failed")
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessID, ok := claims["access_id"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessID: accessID,
			UserID:   userID,
		}, nil
	}
	return nil, err
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {

	access, ok := ExtractToken(r)
	if !ok {
		return nil, errors.New("Can't find access token from Bearer Authorization header")
	}

	cfg := config.Cfg.Jwt
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtAccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m *Token) Response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	resp := NewRespMsg(0, "OK", struct {
		Access           string `json:"access_token"`
		TokenType        string `json:"token_type"`
		AccessExpiresIn  string `json:"access_expires_in"`
		Refresh          string `json:"refresh_token"`
		RefreshExpiresIn string `json:"refresh_expires_in"`
	}{
		m.Access,
		"Bearer",
		fmt.Sprintf("%.0f", m.AccessExpiresIn.Seconds()),
		m.Refresh,
		fmt.Sprintf("%.0f", m.RefreshExpiresIn.Seconds()),
	})

	resp.WriteTo(w)
}
