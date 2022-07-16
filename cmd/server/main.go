package main

import (
	"github.com/gin-gonic/gin"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
	"github.com/scarlettmiss/engine-w/socket"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/home.html")
}

func main() {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	app, err := application.New(sessionRepo, userRepo)
	if err != nil {
		panic(err)
	}

	//create websocket server
	wsAPI, err := socket.New(app)
	if err != nil {
		panic(err)
	}

	wsAPI.CreateHandlers()

	router.GET("/socket.io/*any", gin.WrapH(wsAPI.Server))
	router.POST("/socket.io/*any", gin.WrapH(wsAPI.Server))
	router.GET("/", func(ctx *gin.Context) {
		serveHome(ctx.Writer, ctx.Request)
	})

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	// Start listening and serving requests
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}

	waitForInterrupt := make(chan os.Signal, 1)
	signal.Notify(waitForInterrupt, os.Interrupt, os.Kill)

	<-waitForInterrupt
	defer wsAPI.Server.Close()
}
