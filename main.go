package main

import (
	"fmt"
	"slang/server"
)
//var client  *http.Client
//func main() {
//	client = &http.Client{Timeout: 10 * time.Second}
//	api.FindActivitiesFromApi()
//}

func main() {
	serve := &server.Server{}
	fmt.Println("starting server")
	serve.Start()

}