package handler

import (
	"context"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"jwtToken/service/userService"
	"jwtToken/util"
	"net/http"
)

// SignInHandler : 注册用户接口
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//获取参数
	user := &userService.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email: sql.NullString{
			r.FormValue("email"),
			true,
		},
		Phone: sql.NullString{
			r.FormValue("phone"),
			true,
		},
	}
	ctx := context.Background()
	// 将用户信息注册到用户表中
	if err := userService.Service.Create(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	resp := util.NewRespMsg(0, "OK", "signup succeed")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteTo(w)
}
func SignoutHandler(w http.ResponseWriter, r *http.Request) {

	acessD, err := userService.Service.TokenService.ExtractTokenMetadata(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	err = userService.Service.TokenService.DeleteTokens(acessD)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//respone
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := util.NewRespMsg(0, "OK", "Successfully signout")
	resp.WriteTo(w)
}
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//获取参数
	userName := r.FormValue("username")
	userPass := r.FormValue("password")

	ctx := context.Background()

	//验证密码
	err := userService.Service.Verify(ctx, userName, userPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := userService.Service.GetByName(ctx, userName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//生成token
	token, err := userService.Service.TokenService.NewToken(user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = userService.Service.TokenService.CacheAuth(user.UserID, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	//respone
	token.Response(w)

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
	token, err := userService.Service.TokenService.ParseRefreshToken(refresh)
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
		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "invalid Refresh token", http.StatusUnprocessableEntity)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := userService.Service.TokenService.DeleteCacheAuth(refreshID)
		if delErr != nil || deleted == 0 { //可能是之前登出或取消授权，所以禁止刷新
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		//Create new pairs of refresh and access tokens
		token, err := userService.Service.TokenService.NewToken(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		//save the tokens metadata to redis
		err = userService.Service.TokenService.CacheAuth(userID, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		//respone
		token.Response(w)
	} else {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
	}

}
