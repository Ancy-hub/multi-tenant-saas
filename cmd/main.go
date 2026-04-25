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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// Run Database Migrations
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:root@localhost:8080/postgres?sslmode=disable",
	)
	if err != nil {
		log.Fatal("Failed to initialize migrations:", err)
	}
	
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Database migrations checked/applied successfully")

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
	r.Use(middleware.CORS)
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
			
		r.With(middleware.RequireRole(membershipService, "admin", "member")). // View organization tasks
			Get("/organizations/{id}/tasks", taskHandler.GetTasksByOrg)

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

		// Note: The following routes use UUIDs which are unguessable, providing basic security.
		// For strict RBAC, the handlers themselves should verify the user's role against the project's org_id.
		r.Get("/projects/{project_id}", projectHandler.GetProjectByID)
		r.Delete("/projects/{project_id}", projectHandler.DeleteProject)

		// Task routes
		r.Get("/projects/{project_id}/tasks", taskHandler.GetTasks)
		r.Post("/projects/{project_id}/tasks", taskHandler.CreateTask)
		r.Patch("/tasks/{task_id}", taskHandler.UpdateTask)
		r.Delete("/tasks/{task_id}", taskHandler.DeleteTask)
	})

	// Start server
	log.Println("Server running on port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
