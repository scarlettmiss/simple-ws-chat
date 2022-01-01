package main

import (
	"github.com/scarlettmiss/engine-w/application"
)

func main() {
	app, err := application.New(nil, nil)
	if err != nil {
		panic(err)
	}
	if app == nil {
		panic("app is nil")
	}
}
