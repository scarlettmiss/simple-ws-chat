package main

import (
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
)

func main() {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	app, err := application.New(sessionRepo, userRepo)
	if err != nil {
		panic(err)
	}
	if app == nil {
		panic("app is nil")
	}
}
