package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"jwtToken/cfg"
	"jwtToken/db"
	"jwtToken/util"
	"net/http"
	"strconv"
)

// SignInHandler : 注册用户接口
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//获取参数
	user := &User{}
	//json.NewDecoder(r.Body).Decode(&user)
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	if valid := user.IsValid(); !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 对密码进行bcrypt加密
	user.Password, err = util.EncodePass(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 将用户信息注册到用户表中
	if err := db.NewUser(user.Username, user.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	resp := util.NewRespMsg(0, "OK", "signup succeed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = resp.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func SignoutHandler(w http.ResponseWriter, r *http.Request) {

	acessD, err := util.ExtractTokenMetadata(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	err = util.DeleteTokens(acessD)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//respone
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//resp := util.NewRespMsg(0, "OK", "user "+token.Claims.(*util.Claims).Username+"signout succeed")
	resp := util.NewRespMsg(0, "OK", "Successfully signout")
	if _, err := resp.WriteTo(w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//获取参数
	user := &User{}
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	//验证密码
	if suc := user.Verify(); !suc {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//生成token
	token, err := util.NewToken(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = token.CacheAuth(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	//respone
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := util.NewRespMsg(0, "OK", token)
	_, err = resp.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello,welcome!"))
	return
}

//RefreshTokenHandler :刷新token
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refresh := r.FormValue("refresh-token")
	if refresh == "" {
		http.Error(w, "Can't found refresh token", http.StatusBadRequest)
		return
	}

	//verify the token
	cfg := cfg.Cfg.Jwt
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtRefreshSecret), nil
	})

	if err != nil {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		http.Error(w, "invalid Refresh token", http.StatusUnauthorized)
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshID, ok := claims["refresh_id"].(string) //convert the interface to string
		if !ok {
			http.Error(w, "invalid Refresh token", http.StatusUnprocessableEntity)
			return
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			http.Error(w, "invalid Refresh token", http.StatusUnprocessableEntity)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := util.DeleteCacheAuth(refreshID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		//Create new pairs of refresh and access tokens
		token, err := util.NewToken(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		//save the tokens metadata to redis
		err = token.CacheAuth(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		//respone
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp := util.NewRespMsg(0, "OK", token)
		_, err = resp.WriteTo(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
	}

}
