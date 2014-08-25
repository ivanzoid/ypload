package ylogin

import (
	"fmt"
	"github.com/braintree/manners"
	"net/http"
)

type TokenData struct {
	token string
}

type OauthHandler struct {
	server    *manners.GracefulServer
	tokenData *TokenData
}

func (handler OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.tokenData.token = r.URL.RawQuery
	w.Write([]byte("OK"))
	handler.server.Shutdown <- true
}

func acquireOauthToken(localHttpServerPort int, tokenChan chan string) {
	server := manners.NewServer()
	addressString := fmt.Sprintf(":%d", localHttpServerPort)
	handler := OauthHandler{server: server, tokenData: &TokenData{}}

	server.ListenAndServe(addressString, handler)

	fmt.Println("server did shutdown")
	fmt.Println("handler.token = " + handler.tokenData.token)

	tokenChan <- handler.tokenData.token
}

func Login(localHttpServerPort int, tokenChan chan string) {
	go acquireOauthToken(localHttpServerPort, tokenChan)
}
