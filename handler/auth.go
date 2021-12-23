package handler

import (
	"jwtToken/util"
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := util.ExtractTokenMetadata(r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			//token校验通过: 请求handler处理
			h(w, r)
		})
}
