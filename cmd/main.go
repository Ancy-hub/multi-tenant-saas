package main

import (
	"log"
	"net/http"

	"github.com/ancy-shibu/multi-tenant-saas/internal/config"
	"github.com/ancy-shibu/multi-tenant-saas/internal/db"
	"github.com/ancy-shibu/multi-tenant-saas/internal/handlers"
	"github.com/ancy-shibu/multi-tenant-saas/internal/middleware"
	"github.com/ancy-shibu/multi-tenant-saas/internal/repository"
	"github.com/ancy-shibu/multi-tenant-saas/internal/services"
	"github.com/ancy-shibu/multi-tenant-saas/internal/utils"
	"github.com/go-chi/chi/v5"
)

func main() {
	//Load config
	cfg := config.Load()
	utils.SetJWTSecret(cfg.JWTSecret)
	database, err := db.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}
	defer database.Close()
	log.Println("Database connected")

	orgRepo := repository.NewOrganizationRepository(database) //Initialize repository
	orgService := services.NewOrganizationService(orgRepo)    //Initialize service
	orgHandler := handlers.NewOrganizationHandler(orgService) //Initialize handler

	userRepo := repository.NewUserRepository(database)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	membershipRepo := repository.NewMembershipRepository(database)
	membershipService := services.NewMembershipService(membershipRepo)
	membershipHandler := handlers.NewMembershipHandler(membershipService)

	r := chi.NewRouter()                                            //Router
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { //Health check
		w.Write([]byte("OK"))
	})

	r.Post("/login", userHandler.Login)

	//Organization routes
	r.Post("/organizations", orgHandler.CreateOrganization)
	r.Get("/organizations", orgHandler.GetOrganizations)
	r.Get("/organizations/{id}", orgHandler.GetOrganizationByID)
	r.Patch("/organizations/{id}", orgHandler.UpdateOrganization)
	//User routes
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetAllUsers)
	r.Get("/users/{id}", userHandler.GetUserByID)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.LoggingMiddleware)
		r.Use(middleware.RateLimitMiddleware)
		r.With(middleware.RequireRole(membershipService, "admin", "member")). // View (admin + member)
			Get("/organizations/{id}/members", membershipHandler.GetMembersByOrg)

		r.With(middleware.RequireRole(membershipService, "admin")). // Admin-only actions
			Post("/organizations/{id}/members", membershipHandler.AddUser)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Delete("/organizations/{org_id}/members/{user_id}", membershipHandler.RemoveMember)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Patch("/organizations/{org_id}/members/{user_id}", membershipHandler.UpdateRole)
	})
	r.Get("/users/{id}/organizations", membershipHandler.GetUserOrgs)

	log.Println("Server running on port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
