package main

import (
	"fmt"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
	"os"
	"os/signal"
)

func main() {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()

	app, err := application.New(sessionRepo, userRepo)
	if err != nil {
		panic(err)
	}

	//create websocket server
	fmt.Println(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
}
