package main

import (
	"fmt"
	"jwtToken/cache/redis"
	"jwtToken/config"
	"jwtToken/db/mariadb"
	"jwtToken/handler"
	"jwtToken/service/tokenService"
	"jwtToken/service/userService"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.InitCfg("./config/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("connected to redis")

	database, err := mariadb.NewDatabase(cfg.Db)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("connected to database")

	redisStorage := redis.NewRedisStorage(redisClient)

	tokenService := tokenService.NewService(redisStorage, cfg.Jwt)
	userStorage := mariadb.NewUserStorage(database)
	userService.NewService(userStorage, tokenService)

	http.HandleFunc("/signin", handler.SigninHandler)
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/signout", handler.SignoutHandler)
	http.HandleFunc("/refresh", handler.RefreshTokenHandler)
	http.HandleFunc("/welcome", handler.HTTPInterceptor(handler.WelcomeHandler))

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), nil))
}
