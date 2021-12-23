package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"jwtToken/cache/redis"
	"jwtToken/cfg"
	"net/http"
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

type JWTAccessClaims struct {
	jwt.StandardClaims
}

type Token struct {
	Access          string
	AccessCreateAt  time.Time
	AccessExpiresIn time.Duration

	Refresh          string
	RefreshCreateAt  time.Time
	RefreshExpiresIn time.Duration
}

//NewToken: 创建Token
func NewToken(username string) (*Token, error) {
	var err error
	now := time.Now()
	cfg := cfg.Cfg.Jwt
	accessDuration, _ := time.ParseDuration(cfg.AccessTokenDuration)
	refreshDuration, _ := time.ParseDuration(cfg.RefreshTokenDuration)
	token := &Token{
		AccessCreateAt:   now,
		AccessExpiresIn:  accessDuration,
		RefreshCreateAt:  now,
		RefreshExpiresIn: refreshDuration,
	}

	// Create the access JWT claims
	claims := &JWTAccessClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: username,
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: token.AccessCreateAt.Add(token.AccessExpiresIn).Unix(),
		},
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JwtAccessSecret))
	if err != nil {
		return nil, err
	}
	token.Access = access

	// Create the refresh JWT claims
	claims = &JWTAccessClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: username,
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: token.RefreshCreateAt.Add(token.RefreshExpiresIn).Unix(),
		},
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JwtRefreshSecret))
	if err != nil {
		return nil, err
	}
	token.Refresh = refresh

	return token, nil
}

func (m *Token) SetCookieAuth(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    m.Access,
		HttpOnly: true, //防止 XSS attack
		Expires:  m.AccessCreateAt.Add(m.AccessExpiresIn),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    m.Refresh,
		HttpOnly: true, //防止 XSS attack
		Expires:  m.RefreshCreateAt.Add(m.RefreshExpiresIn),
	})
}

func UnsetCookieAuth(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    "",
		HttpOnly: true, //防止 XSS attack
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		HttpOnly: true, //防止 XSS attack
		MaxAge:   -1,
	})
}

//获取accessToken顺序：
//1.Authorization头
//2.formValue中access-token
//3.cookie中access-token
func BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = r.FormValue("access-token")
	}

	if token == "" {
		c, err := r.Cookie("access-token")

		if err != nil {
			return "", false
		}

		return c.Value, true

	}
	return token, token != ""
}

func (m *Token) CacheAuth() error {
	var err error

	username, err := ExtractTokenSubject(m.Access)
	if err != nil {
		return err
	}

	err = redis.RedisClient().Set("access-"+username, m.Access, m.AccessExpiresIn).Err()
	if err != nil {
		return err
	}
	err = redis.RedisClient().Set("refresh-"+username, m.Refresh, m.RefreshExpiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

func DelCacheAuth(username string) {
	redis.RedisClient().Del("access-" + username)
	redis.RedisClient().Del("refresh-" + username)
}

func ExtractTokenSubject(token string) (string, error) {
	jwtToken, err := VerifyToken(token, false)
	if err != nil {
		return "", err
	}

	return jwtToken.Claims.(*JWTAccessClaims).StandardClaims.Subject, nil

}
func VerifyToken(token string, isRefresh bool) (*jwt.Token, error) {

	claims := &JWTAccessClaims{}
	cfg := cfg.Cfg.Jwt
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if !isRefresh {
			return []byte(cfg.JwtAccessSecret), nil
		}
		return []byte(cfg.JwtRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		if !isRefresh {
			return nil, errors.New("invalid token")
		}
		return nil, errors.New("invalid refresh token")
	}
	return tkn, nil
}
