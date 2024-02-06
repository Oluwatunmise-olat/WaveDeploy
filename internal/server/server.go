package server

import (
	"errors"
	"fmt"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/http_server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func BootStrapHttpServer() {
	router := mux.NewRouter().StrictSlash(false)

	http_server.AddRouteHandlers(router)
	_server := &http.Server{Addr: fmt.Sprintf(":%s", os.Getenv("PORT")), Handler: router}

	log.Println("🦍Http server started 🦍")
	if err := _server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("An error occurred starting http_server server. Error: %s", err)
	}

}
