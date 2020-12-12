package main

import (
	"net/http"

	"google.golang.org/api/oauth2/v2"
)

var (
	httpClient    = &http.Client{}
	googleOauthId = ""
)

type Auth struct {
	Token string `json:"token"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		token := req.Header.Get("Authorization")
		// tokenResponse, error := verifyIdToken(token)
		// if error == nil && tokenResponse.Audience == googleOauthId {
		// 	next.ServeHTTP(w, req)
		// } else {
		// 	sendError(w, http.StatusUnauthorized, "No token found")
		// }
		if token == "" {
			sendError(w, http.StatusUnauthorized, "No token found")
		} else {
			next.ServeHTTP(w, req)
		}
	})
}

func verifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
