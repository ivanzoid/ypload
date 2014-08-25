package ylogin

import (
	"fmt"
	"github.com/braintree/manners"
	"net/http"
)

const (
	kRedirectHtmlTemplate = "<script> window.location = document.URL.split('#').join('?') </script>"
	kDoneHtmlTemplate     = `<style type="text/css">
								html,
								body {
									width: 100%;
									height: 100%;
								}
								html {
									display: table;
								}
								body {
									display: table-cell;
									vertical-align: middle;
								}
							</style>
							<body>
							<center>Great, we've got a token! You now may return back to terminal.</center>
							</body>`
	kTokenKey     = "access_token"
	kExpiresInKey = "expires_in"
)

type TokenData struct {
	Token     string
	ExpiresIn string
}

type OauthHandler struct {
	server    *manners.GracefulServer
	tokenData *TokenData
	doneChan  chan bool
}

func (handler OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if len(params) == 0 {
		fmt.Printf("Redirecting...\n")
		fmt.Fprintf(w, "%v", kRedirectHtmlTemplate)
	} else {
		handler.tokenData.Token = params.Get(kTokenKey)
		handler.tokenData.ExpiresIn = params.Get(kExpiresInKey)
		fmt.Fprintf(w, "%v", kDoneHtmlTemplate)
		fmt.Printf("Done\n")
		handler.server.Shutdown <- true
		handler.doneChan <- true
	}
}

func acquireOauthToken(localHttpServerPort int, tokenDataChan chan TokenData) {
	server := manners.NewServer()
	addressString := fmt.Sprintf(":%d", localHttpServerPort)
	handler := OauthHandler{server: server, tokenData: &TokenData{}, doneChan: make(chan bool)}

	go server.ListenAndServe(addressString, handler)

	<-handler.doneChan

	fmt.Println("TokenData: %v", handler.tokenData)

	tokenDataChan <- *handler.tokenData
}

func Login(localHttpServerPort int, tokenDataChan chan TokenData) {
	go acquireOauthToken(localHttpServerPort, tokenDataChan)
}
