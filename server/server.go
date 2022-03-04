package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "os"
)

type Server struct {
	Router *gin.Engine
}

func defineRoutes(router *gin.Engine) {
	apiRouter := router.Group("/api/v1")
	{
		apiRouter.GET("/sessions", GetActivities())

	}
}

func (server Server) Start() {
	r := gin.New()

	defineRoutes(r)

	PORT := fmt.Sprintf(":%s", "8080")
	if PORT == ":" {
		PORT = ":8080"
	}
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Printf("Server started on %s\n", PORT)
}