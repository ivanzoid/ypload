package ylogin

import (
	"fmt"
	"github.com/braintree/manners"
	"net/http"
)

const (
	kLocalHttpServerPort = 30171
)

type TokenData struct {
	token string
}

type OauthHandler struct {
	server    *manners.GracefulServer
	tokenData *TokenData
}

func (handler OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("query = %v\n", r.URL.RawQuery)
	fmt.Printf("fragment = %v\n", r.URL.Fragment)
	handler.tokenData.token = r.URL.RawQuery
	w.Write([]byte("OK"))
	handler.server.Shutdown <- true
}

func acquireOauthToken(tokenChan chan string) {
	server := manners.NewServer()
	addressString := fmt.Sprintf(":%d", kLocalHttpServerPort)
	handler := OauthHandler{server: server, tokenData: &TokenData{}}

	server.ListenAndServe(addressString, handler)

	fmt.Println("server did shutdown")
	fmt.Println("handler.token = " + handler.tokenData.token)

	tokenChan <- handler.tokenData.token
}

func Login(tokenChan chan string) {
	go acquireOauthToken(tokenChan)
}
