package main

import (
	"log"

	"my_law_server/internal/config"
	"my_law_server/internal/database"
	"my_law_server/internal/handler"
	"my_law_server/internal/repository"
	"my_law_server/internal/server"
	"my_law_server/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.OpenReadOnlySQLite(cfg.LawDBPath)
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}

	typeRepo := repository.NewTypeRepository(db)
	lawRepo := repository.NewLawRepository(db)
	parsedLawRepo := repository.NewParsedLawRepository(cfg.LawDetailJSONDir)
	lawService := service.NewLawService(typeRepo, lawRepo, parsedLawRepo)
	lawHandler := handler.NewLawHandler(lawService)
	router := server.NewRouter(lawHandler)

	log.Printf("law server listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run http server: %v", err)
	}
}
