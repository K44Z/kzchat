package configs

import (
	"fmt"
	"os"

	"github.com/K44Z/kzchat/internal/api"
)

func SetBaseUrl() {
	port := os.Getenv("PORT")
	if len(port) > 0 && port[0] == ':' {
		port = port[1:]
	}
	api.BASE_URL = fmt.Sprintf("http://localhost:%s", port)
	api.WS_URL = fmt.Sprintf("ws://localhost:%s", port)
}
