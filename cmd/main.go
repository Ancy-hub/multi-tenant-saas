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

	
	orgRepo:= repository.NewOrganizationRepository(database)//Initialize repository
	orgService:= services.NewOrganizationService(orgRepo)//Initialize service
	orgHandler:= handlers.NewOrganizationHandler(orgService)//Initialize handler
	
	userRepo:=repository.NewUserRepository(database)
	userService:=services.NewUserService(userRepo)
	userHandler:=handlers.NewUserHandler(userService)

	membershipRepo := repository.NewMembershipRepository(database)
	membershipService := services.NewMembershipService(membershipRepo)
	membershipHandler := handlers.NewMembershipHandler(membershipService)
	
	r := chi.NewRouter()//Router
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {//Health check
		w.Write([]byte("OK"))
	})

	//Organization routes
	r.Post("/organizations",orgHandler.CreateOrganization)
	r.Get("/organizations",orgHandler.GetOrganizations)
	r.Get("/organizations/{id}",orgHandler.GetOrganizationByID)
	r.Patch("/organizations/{id}", orgHandler.UpdateOrganization)
	//User routes
	r.Post("/users",userHandler.CreateUser)
	r.Get("/users",userHandler.GetAllUsers)
	r.Get("/users/{id}",userHandler.GetUserByID)

	r.Post("/organizations/{id}/members", membershipHandler.AddUser)
	
	log.Println("Server running on port", cfg.Port)
	err= http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}