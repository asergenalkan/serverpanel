package api

import (
	"github.com/asergenalkan/serverpanel/internal/config"
	"github.com/asergenalkan/serverpanel/internal/database"
	"github.com/asergenalkan/serverpanel/internal/middleware"
	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	db  *database.DB
	cfg *config.Config
}

func SetupRoutes(router fiber.Router, db *database.DB) {
	cfg := config.Load()
	h := &Handler{db: db, cfg: cfg}

	// Public routes
	router.Post("/auth/login", h.Login)
	router.Get("/health", h.Health)

	// Protected routes
	protected := router.Group("/", middleware.AuthMiddleware(cfg.JWTSecret))

	// Auth
	protected.Get("/auth/me", h.GetCurrentUser)
	protected.Post("/auth/logout", h.Logout)

	// Dashboard
	protected.Get("/dashboard/stats", h.GetDashboardStats)

	// Users (admin only)
	adminOnly := protected.Group("/", middleware.RoleMiddleware(models.RoleAdmin))
	adminOnly.Get("/users", h.ListUsers)
	adminOnly.Post("/users", h.CreateUser)
	adminOnly.Get("/users/:id", h.GetUser)
	adminOnly.Put("/users/:id", h.UpdateUser)
	adminOnly.Delete("/users/:id", h.DeleteUser)

	// Packages (admin only)
	adminOnly.Get("/packages", h.ListPackages)
	adminOnly.Post("/packages", h.CreatePackage)
	adminOnly.Put("/packages/:id", h.UpdatePackage)
	adminOnly.Delete("/packages/:id", h.DeletePackage)

	// Accounts - Hosting hesapları (admin only)
	adminOnly.Get("/accounts", h.ListAccounts)
	adminOnly.Post("/accounts", h.CreateAccount)
	adminOnly.Get("/accounts/:id", h.GetAccount)
	adminOnly.Delete("/accounts/:id", h.DeleteAccount)
	adminOnly.Post("/accounts/:id/suspend", h.SuspendAccount)
	adminOnly.Post("/accounts/:id/unsuspend", h.UnsuspendAccount)

	// Domains
	protected.Get("/domains", h.ListDomains)
	protected.Post("/domains", h.CreateDomain)
	protected.Get("/domains/:id", h.GetDomain)
	protected.Delete("/domains/:id", h.DeleteDomain)

	// Databases
	protected.Get("/databases", h.ListDatabases)
	protected.Post("/databases", h.CreateDatabase)
	protected.Delete("/databases/:id", h.DeleteDatabase)

	// System (admin only)
	adminOnly.Get("/system/stats", h.GetSystemStats)
	adminOnly.Get("/system/services", h.GetServices)
	adminOnly.Post("/system/services/:name/restart", h.RestartService)
}
