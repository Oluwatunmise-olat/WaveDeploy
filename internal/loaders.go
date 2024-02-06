package internal

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/joho/godotenv"
)

func BootstrapApp() {
	godotenv.Load()
	_, err := db.Connect()

	if err != nil {
		panic("Error Establishing Database Connection 🚧")
	}
}
