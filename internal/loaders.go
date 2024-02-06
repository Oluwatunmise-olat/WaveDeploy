package internal

import (
	"flag"
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/server"
	"github.com/joho/godotenv"
)

func BootstrapApp() {
	serveHTTP := flag.Bool("serve-http", false, "start HTTP server")
	flag.Parse()

	godotenv.Load()

	_, err := db.Connect()
	if err != nil {
		panic("Error Establishing Database Connection 🚧")
	}

	if *serveHTTP {
		server.BootStrapHttpServer()
		return
	}

	cmd.Execute()
}
