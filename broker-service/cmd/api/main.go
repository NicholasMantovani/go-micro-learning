package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const webPort = "80"

type Config struct {
	AuthBaseUrl string
}

func main() {

	authBaseUrl := os.Getenv("AUTH_BASE_URL")
	app := Config{
		AuthBaseUrl: authBaseUrl,
	}
	log.Printf("Starting broken server on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the serer
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
