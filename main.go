package main

import (
	"fmt"
	"jwtToken/cache/redis"
	"jwtToken/cfg"
	"jwtToken/db/mariadb"
	"jwtToken/handler"
	"log"
	"net/http"
)

func main() {
	cfg.InitCfg("./cfg/cfg.yaml")
	redis.InitRedis()
	mariadb.InitDB()

	// 静态资源处理
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/signin", handler.SigninHandler)
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/signout", handler.SignoutHandler)
	http.HandleFunc("/refresh", handler.RefreshTokenHandler)
	http.HandleFunc("/welcome", handler.HTTPInterceptor(handler.WelcomeHandler))

	cfg := cfg.Cfg.Server

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
