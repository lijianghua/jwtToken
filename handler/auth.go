package handler

import (
	"jwtToken/service/userService"
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authD, err := userService.Service.TokenService.ExtractTokenMetadata(r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			_, err = userService.Service.TokenService.FetchAuth(authD)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			//token校验通过: 请求handler处理
			h(w, r)
		})
}
