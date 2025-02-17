package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"think-before-you-speak/pkg/config"
	server "think-before-you-speak/pkg/http"
	"think-before-you-speak/pkg/log"
)

var (
	printVersion = flag.Bool("v", false, "print version")
	appConfig    = flag.String("config", "../config.yaml", "application config path")
	isLocal      = flag.Bool("local", true, "local mode as compared to when deployed on managed Service")
)

// @title Think Before You Speak Swagger API
// @version 1.0
// @description This is the API documentation for Think Before You Speak

// @contact.name API Support
// @contact.email dbaekajnr@gmail.com

// @license.name MIT License
// @license.url http://opensource.org/licenses/MIT

// @securityDefinitions.apikey OpenRouterApiKey
// @in header
// @name X-OpenRouter-API-Token
// @description API key for DeepSeek

// @securityDefinitions.apikey AnthropicApiKeyAuth
// @in header
// @name X-Anthropic-API-Token
// @description API key for Anthropic model
func main() {
	flag.Parse()

	if *printVersion {
		log.PrintVersion()
		os.Exit(0)
	}

	if *isLocal {
		err := godotenv.Load("../.env")
		if err != nil {
			panic(err)
		}
	}

	l := log.Get()

	conf, err := config.Parse(*appConfig)
	if err != nil {
		l.Fatal().Err(err).Msgf("Failed to parse config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	e, err := server.NewEcho(conf)
	if err != nil {
		l.Fatal().Err(err).Msgf("Failed to create server: %v", err)
	}

	l.Info().Msg("Starting server")

	if err = e.Run(); err != nil {
		l.Fatal().Err(err).Msgf("Failed to run server: %v", err)
	}
}
