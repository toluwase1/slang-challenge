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
	//router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (server Server) Start() {
	r := gin.New()

	defineRoutes(r)

	//PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	PORT := fmt.Sprintf(":%s", "8080")
	if PORT == ":" {
		PORT = ":8080"
	}
	fmt.Println("port here",PORT)
	srv := &http.Server{
		Addr:    PORT,
		Handler: r,
	}
	fmt.Println("key",srv)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("i got here")
		log.Fatalf("listen: %s\n", err)
	}
	log.Printf("Server started on %s\n", PORT)
}