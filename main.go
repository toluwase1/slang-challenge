package main

import (
	"github.com/joho/godotenv"
	"slang/server"
)
//var client  *http.Client
//func main() {
//	client = &http.Client{Timeout: 10 * time.Second}
//	api.FindActivitiesFromApi()
//}

func main() {
	_ = godotenv.Load()
	serve := &server.Server{}
	serve.Start()
}