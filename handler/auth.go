package handler

import (
	"jwtToken/util"
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//get access token
			tokenStr, suc := util.BearerAuth(r)
			if !suc {
				http.Error(w, "Not Authorized", http.StatusUnauthorized)
				return
			}
			//verify token
			//token, err := util.VerifyToken(tokenStr)
			_, err := util.VerifyToken(tokenStr, false)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			//token校验通过: 请求handler处理
			h(w, r)
		})
}
