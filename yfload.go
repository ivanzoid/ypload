package main

import (
	"./config"
	"./ylogin"
	"fmt"
	// "github.com/skratchdot/open-golang/open"
)

func getToken() string {
	tokenChan := make(chan string)
	ylogin.Login(tokenChan)
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
