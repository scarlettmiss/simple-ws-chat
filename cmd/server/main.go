package main

import (
	"flag"
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

var addr = flag.String("addr", ":8080", "http service address")

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

	app, err := application.New(sessionRepo, userRepo)
	if err != nil {
		panic(err)
	}

	//create websocket server
	wsAPI, err := socket.New(app)
	if err != nil {
		panic(err)
	}

	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	// Create API route group
	api := router.Group("/")
	{
		api.GET("/", func(ctx *gin.Context) {
			serveHome(ctx.Writer, ctx.Request)
		})
		// Add /hello GET route to router and define route handler function
		api.GET("/ws", func(ctx *gin.Context) {
			//ctx.JSON(200, gin.H{"msg": "world"})
			wsAPI.Handle(ctx.Writer, ctx.Request)
		})
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	// Start listening and serving requests
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}

	//http.HandleFunc("/", serveHome)
	//http.HandleFunc(
	//	"/ws", func(w http.ResponseWriter, req *http.Request) {
	//		wsAPI.Handle(w, req)
	//	},
	//)
	//err = http.ListenAndServe(*addr, nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
	//fmt.Println(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
}
