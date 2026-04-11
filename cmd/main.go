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

// main initializes and starts the multi-tenant SaaS server.
func main() {
	// Load configuration
	cfg := config.Load()
	utils.SetJWTSecret(cfg.JWTSecret)

	// Initialize database connection
	database, err := db.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}
	defer database.Close()
	log.Println("Database connected")

	// Initialize repositories
	orgRepo := repository.NewOrganizationRepository(database)
	userRepo := repository.NewUserRepository(database)
	membershipRepo := repository.NewMembershipRepository(database)
	projectRepo := repository.NewProjectRepository(database)
	taskRepo := repository.NewTaskRepository(database)

	// Initialize services
	orgService := services.NewOrganizationService(orgRepo)
	userService := services.NewUserService(userRepo)
	membershipService := services.NewMembershipService(membershipRepo)
	projectService := services.NewProjectService(projectRepo)
	taskService := services.NewTaskService(taskRepo)

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)
	userHandler := handlers.NewUserHandler(userService)
	membershipHandler := handlers.NewMembershipHandler(membershipService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Initialize router
	r := chi.NewRouter()

	// Health check endpoint (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Authentication endpoint (no auth required)
	r.Post("/login", userHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.LoggingMiddleware)
		r.Use(middleware.RateLimitMiddleware)

		// User routes
		r.Post("/users", userHandler.CreateUser)
		r.Get("/users", userHandler.GetAllUsers)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Get("/users/{id}/organizations", membershipHandler.GetUserOrgs)

		// Organization routes
		r.Post("/organizations", orgHandler.CreateOrganization)
		r.Get("/organizations", orgHandler.GetOrganizations)
		r.Get("/organizations/{id}", orgHandler.GetOrganizationByID)
		r.Patch("/organizations/{id}", orgHandler.UpdateOrganization)

		// Membership routes (organization members)
		r.With(middleware.RequireRole(membershipService, "admin", "member")). // View (admin + member)
			Get("/organizations/{id}/members", membershipHandler.GetMembersByOrg)

		r.With(middleware.RequireRole(membershipService, "admin")). // Admin-only actions
			Post("/organizations/{id}/members", membershipHandler.AddUser)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Delete("/organizations/{org_id}/members/{user_id}", membershipHandler.RemoveMember)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Patch("/organizations/{org_id}/members/{user_id}", membershipHandler.UpdateRole)

		// Project routes
		r.With(middleware.RequireRole(membershipService, "admin", "member")). // View projects (admin + member)
			Get("/organizations/{id}/projects", projectHandler.GetProjects)

		r.With(middleware.RequireRole(membershipService, "admin")). // Admin only
			Post("/organizations/{id}/projects", projectHandler.CreateProject)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Delete("/projects/{project_id}", projectHandler.DeleteProject)

		// Task routes
		r.With(middleware.RequireRole(membershipService, "admin", "member")).
			Get("/projects/{project_id}/tasks", taskHandler.GetTasks)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Post("/projects/{project_id}/tasks", taskHandler.CreateTask)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Patch("/tasks/{task_id}", taskHandler.UpdateTask)

		r.With(middleware.RequireRole(membershipService, "admin")).
			Delete("/tasks/{task_id}", taskHandler.DeleteTask)
	})

	// Start server
	log.Println("Server running on port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
