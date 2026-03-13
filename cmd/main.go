package main

import (
	"log"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/config"
	"github.com/ancy-shibu/multi-tenant-saas/internal/db"
	"github.com/ancy-shibu/multi-tenant-saas/internal/handlers"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/go-chi/chi/v5"
)

func main() {
	//Load config
	cfg:=config.Load()

	database,err:=db.NewDB();
	if err!=nil{
		log.Fatal("Failed to connect to db:",err)
	}
	defer database.Close()
	log.Println("Database connected")

	//Initialize repository
	orgRepo:= repository.NewOrganizationRepository(database)

	//Initialize service
	orgService:= services.NewOrganizationService(orgRepo)

	//Initialize handler
	orgHandler:= handlers.NewOrganizationHandler(orgService)

	//Router
	r := chi.NewRouter()

	//Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	//Organization routes
	r.Post("/organizations",orgHandler.CreateOrganization)

	r.Get("/organizations",orgHandler.GetOrganizations)

	r.Get("/organizations/{id}",orgHandler.GetOrganizationByID)

	r.Patch("/organizations/{id}", orgHandler.UpdateOrganization)
	
	log.Println("Server running on port", cfg.Port)

	err= http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}