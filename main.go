package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	GiteaURL      string `env:"GITEA_URL,required"`
	GiteaUser     string `env:"GITEA_USER,required"`
	GiteaPassword string `env:"GITEA_PWD,required"`
	Repo          string `env:"GIT_REPO,required"`
	PullRequestID int    `env:"PULL_REQUEST,required"`
	HTTPTimeout   int    `env:"HTTP_TIMEOUT" envDefault:"30"`
}

var cfg config

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	gitea := Gitea{
		BaseURL: cfg.GiteaURL,
		Client: &http.Client{
			Timeout: time.Duration(cfg.HTTPTimeout) * time.Second,
		},
		Username: cfg.GiteaUser,
		Password: cfg.GiteaPassword,
	}

	issues := ReadIssues(os.Stdin)
	log.Printf("Found %d issues\n", len(issues))
	review := FormatReview(issues)

	if err := gitea.DiscardPreviousReviews(cfg.Repo, cfg.PullRequestID); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err := gitea.SendReview(cfg.Repo, cfg.PullRequestID, review); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if len(issues) > 0 {
		os.Exit(1)
	}
}