package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const webPort = "80"

type Config struct {
	AuthBaseUrl   string
	LoggerBaseUrl string
	MailerBaseUrl string
}

func main() {

	authBaseUrl := os.Getenv("AUTH_BASE_URL")
	loggerBaseUrl := os.Getenv("LOGGER_BASE_URL")
	mailerBaseUrl := os.Getenv("MAILER_BASE_URL")
	app := Config{
		AuthBaseUrl:   authBaseUrl,
		LoggerBaseUrl: loggerBaseUrl,
		MailerBaseUrl: mailerBaseUrl,
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
