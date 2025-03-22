package main

import (
	"fishbase/routers"
	"log"
)

func main() {
	r := routers.SetupRouter()
	log.Printf("Starting server on port %s", routers.ServerPort)
	log.Printf("Open http://localhost%s in your browser", routers.ServerPort)
	if err := r.Run(routers.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
