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

func getTokenData() ylogin.TokenData {
	tokenDataChan := make(chan ylogin.TokenData)
	ylogin.Login(kLocalHttpServerPort, tokenDataChan)

	openLoginPage(kAppId)

	tokenData := <-tokenDataChan

	return tokenData
}

func main() {
	cfg, _ := config.Load()

	var token string

	if cfg == nil {
		fmt.Println("no oauth token")
		tokenData := getTokenData()
		token = tokenData.Token
	} else {
		token = cfg.OauthToken
	}

	fmt.Printf("token: %v\n", token)
}
