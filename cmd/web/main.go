package main

import (
	"flag"
	"html/template"
	"log"
	"os"

	"github.com/ADEMOLA200/ECommerce/cmd/models"
	"github.com/ADEMOLA200/ECommerce/cmd/server"
)

const (
	version    = "1.0.0"
	cssVersion = "1"
)

func main() {
	var cfg models.Config

	flag.IntVar(&cfg.Port, "port", 4000, "Listening to server...")
	flag.StringVar(&cfg.Env, "env", "development", "Application envionment {development|production}")
	flag.StringVar(&cfg.API, "api", "http://localhost:4001", "URL to api}")

	flag.Parse()

	cfg.Stripe.Key = os.Getenv("STRIPE_KEY")
	cfg.Stripe.Secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc := make(map[string]*template.Template)

	app := &models.Application{
		Config:        cfg,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		TemplateCache: tc,
		Version:       version,
	}

	srv := server.Server{
		App: app,
	}

	srv.CheckAndKillProcess(cfg.Port)

	// Start Server
	err := srv.Server()
	if err != nil {
		errorLog.Println(err)
		errorLog.Fatal(err)
	}
}
