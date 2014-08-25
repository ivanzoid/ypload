package main

import (
	"./config"
	"./ylogin"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
)

const (
	kAppId               = "e2b26273dab84121bf3f9c2be4bb8915"
	kLocalHttpServerPort = 30171
)

func openLoginPage(appId string) {
	urlString := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=token&client_id=%v", appId)
	fmt.Println("Url: " + urlString)
	err := open.Start(urlString)
	if err != nil {
		log.Fatalf("Can't open browser: %v", err)
	}
}

func getToken() string {
	tokenChan := make(chan string)
	ylogin.Login(kLocalHttpServerPort, tokenChan)

	openLoginPage(kAppId)

	token := <-tokenChan

	return token
}

func main() {
	cfg, _ := config.Load()

	var token string

	if cfg == nil {
		fmt.Println("no oauth token")
		token = getToken()
	} else {
		token = cfg.OauthToken
	}

	fmt.Printf("token: %v\n", token)
}
