package main

import (
	"curriculum/internal/handler"
	"curriculum/internal/router"
	"curriculum/internal/service"
	"curriculum/internal/store"
	"log"
)

func main() {
	dataDir := "./data"
	s, err := store.NewStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	svc := service.NewCurriculumService(s)
	h := handler.NewHandler(s, svc)

	r := router.SetupRouter(h)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
