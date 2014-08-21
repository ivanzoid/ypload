package main

import (
	"./config"
	"fmt"
	// "github.com/skratchdot/open-golang/open"
)

func getOauthToken() {

}

func main() {
	cfg, _ := config.ConfigLoad()
	if cfg == nil {
		fmt.Println("no oauth token")
	}
}
