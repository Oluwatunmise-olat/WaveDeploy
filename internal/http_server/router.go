package http_server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AddRouteHandlers(router *mux.Router) {
	router.HandleFunc("/oauth/github/callback", HandleOauthCallbackWebhook).Methods(http.MethodGet)
}
