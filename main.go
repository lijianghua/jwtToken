package main

import (
	"fmt"
	"jwtToken/cache/redis"
	"jwtToken/config"
	"jwtToken/db/mariadb"
	"jwtToken/handler"
	"log"
	"net/http"
)

func main() {
	config.InitCfg("./config/config.yaml")

	if err := redis.NewClient(&config.Cfg); err != nil {
		log.Fatalln(err)
	}

	if err := mariadb.NewDatabase(&config.Cfg); err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/signin", handler.SigninHandler)
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/signout", handler.SignoutHandler)
	http.HandleFunc("/refresh", handler.RefreshTokenHandler)
	http.HandleFunc("/welcome", handler.HTTPInterceptor(handler.WelcomeHandler))

	cfg := config.Cfg.Server

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
