package main

import (
	"errors"
	"log"
	"os"

	"LawHelperServer/internal/config"
	"LawHelperServer/internal/database"
	"LawHelperServer/internal/handler"
	"LawHelperServer/internal/repository"
	"LawHelperServer/internal/server"
	"LawHelperServer/internal/service"
)

func main() {
	cfg := config.Load()
	logParsedLawDirStatus(cfg.LawDetailJSONDir)

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

	log.Printf("parsed law json dir: %s", cfg.LawDetailJSONDir)
	log.Printf("law server listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run http server: %v", err)
	}
}

func logParsedLawDirStatus(dir string) {
	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			log.Printf("parsed law json path exists but is not a directory: %s", dir)
		}
		return
	}

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("parsed law json dir does not exist yet: %s", dir)
		log.Printf("place parsed files as <versionId>.json under that directory")
		return
	}

	log.Printf("check parsed law json dir %s: %v", dir, err)
}
