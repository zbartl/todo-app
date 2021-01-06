package main

import (
	"fmt"
	"github.com/gobridge-kr/todo-app/server"
	"github.com/gobridge-kr/todo-app/server/database"
	"github.com/gobridge-kr/todo-app/server/utils"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	viper.AddConfigPath("env")
	viper.SetConfigName("dev-env")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	port, baseUrl, db := configureDb()
	jwt := configureJwt()

	mux := http.NewServeMux()
	s := server.New(baseUrl)
	s.ConfigureRoutes(mux, db, jwt)
	s.Serve(mux, port)
}

func configureJwt() *jwtea.Provider {
	jwtConfig := &jwtea.Configuration{
		ThirdPartyConfig: jwtea.ThirdPartyConfiguration{
			Url:   viper.GetString("jwt.third_party.url"),
			ClientId:   viper.GetString("jwt.third_party.cid"),
			ClientSecret:   viper.GetString("jwt.third_party.secret"),
			ThirdPartyAudience:   viper.GetString("jwt.third_party.audience"),
		},
		Secret:   viper.GetString("jwt.secret"),
		Issuer:   viper.GetString("jwt.issuer"),
		Audience: viper.GetString("jwt.audience"),
	}
	return jwtea.NewProvider(jwtConfig)
}

func configureDb() (string, string, *database.Database) {
	port := viper.GetString("todo.port")
	if len(port) == 0 {
		log.Panicf("Error reading config value for port.")
	}

	baseUrl := viper.GetString("todo.baseurl")
	if len(baseUrl) == 0 {
		log.Panicf("Error reading config value for base url.")
	}

	dbConfig := database.Config{
		BaseURL: fmt.Sprintf("%s:%s", baseUrl, port),
	}
	db := database.New(dbConfig)
	return port, baseUrl, db
}
