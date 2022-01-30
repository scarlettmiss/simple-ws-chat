package main

import (
	"fmt"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
	"github.com/scarlettmiss/engine-w/socket"
	"net/http"
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
	wsAPI, err := socket.New(app)
	if err != nil {
		panic(err)
	}
	http.HandleFunc(
		"/ws", func(w http.ResponseWriter, req *http.Request) {
			wsAPI.Handle(w, req)
		},
	)
	fmt.Println(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
}
