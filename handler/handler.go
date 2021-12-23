package handler

import (
	"io/ioutil"
	"jwtToken/db"
	"jwtToken/util"
	"net/http"
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
	//get access token
	access, suc := util.BearerAuth(r)
	if !suc {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	token, err := util.VerifyToken(access, false)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	username := token.Claims.(*util.JWTAccessClaims).StandardClaims.Subject
	//unset auth cookie
	util.UnsetCookieAuth(w)

	util.DelCacheAuth(username)

	//respone
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//resp := util.NewRespMsg(0, "OK", "user "+token.Claims.(*util.Claims).Username+"signout succeed")
	resp := util.NewRespMsg(0, "OK", "signout succeed")
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
	token, err := util.NewToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set auth cookie
	token.SetCookieAuth(w)
	token.CacheAuth()

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

	data, err := ioutil.ReadFile("./static/view/main.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

//RefreshTokenHandler :刷新token
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	//get access token
	access, suc := util.BearerAuth(r)
	if suc {
		//verify token,access token有效则不刷新
		//token, err := util.VerifyToken(tokenStr)
		_, err := util.VerifyToken(access, false)
		if err == nil {
			//respone
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			resp := util.NewRespMsg(0, "OK", "access token still valid! not need refresh")
			_, err = resp.WriteTo(w)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}

	//查找refresh token

	refresh := r.FormValue("refresh-token")
	if refresh == "" {
		c, err := r.Cookie("refresh-token")

		if err == nil {
			refresh = c.Value
		}
	}

	refreshJwt, err := util.VerifyToken(refresh, true)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	username := refreshJwt.Claims.(*util.JWTAccessClaims).StandardClaims.Subject

	//TODO： 增加查询redis中 refreshToken是否存在（登出），不应只检查refreshToken有效无效。
	//此外！ 为减少refreshToken在内存中占用，

	//生成token
	token, err := util.NewToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set auth cookie
	token.SetCookieAuth(w)
	token.CacheAuth()

	//respone
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := util.NewRespMsg(0, "OK", token)
	_, err = resp.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
