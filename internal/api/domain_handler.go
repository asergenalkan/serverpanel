package api

import (
	"strconv"

	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ListDomains(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	role := c.Locals("role").(string)

	var query string
	var args []interface{}

	if role == models.RoleAdmin {
		query = `SELECT id, user_id, name, document_root, ssl_enabled, ssl_expiry, active, created_at FROM domains ORDER BY name`
	} else {
		query = `SELECT id, user_id, name, document_root, ssl_enabled, ssl_expiry, active, created_at FROM domains WHERE user_id = ? ORDER BY name`
		args = append(args, userID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch domains",
		})
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var d models.Domain
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.DocumentRoot, &d.SSLEnabled, &d.SSLExpiry, &d.Active, &d.CreatedAt); err != nil {
			continue
		}
		domains = append(domains, d)
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    domains,
	})
}

func (h *Handler) CreateDomain(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req struct {
		Name         string `json:"name"`
		DocumentRoot string `json:"document_root"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Domain name is required",
		})
	}

	if req.DocumentRoot == "" {
		req.DocumentRoot = "/home/" + c.Locals("username").(string) + "/public_html/" + req.Name
	}

	result, err := h.db.Exec(`
		INSERT INTO domains (user_id, name, document_root, active)
		VALUES (?, ?, ?, 1)
	`, userID, req.Name, req.DocumentRoot)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Domain already exists",
		})
	}

	id, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Domain created successfully",
		Data:    map[string]int64{"id": id},
	})
}

func (h *Handler) GetDomain(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid domain ID",
		})
	}

	var domain models.Domain
	err = h.db.QueryRow(`
		SELECT id, user_id, name, document_root, ssl_enabled, ssl_expiry, active, created_at
		FROM domains WHERE id = ?
	`, id).Scan(&domain.ID, &domain.UserID, &domain.Name, &domain.DocumentRoot, &domain.SSLEnabled, &domain.SSLExpiry, &domain.Active, &domain.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   "Domain not found",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    domain,
	})
}

func (h *Handler) DeleteDomain(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid domain ID",
		})
	}

	userID := c.Locals("user_id").(int64)
	role := c.Locals("role").(string)

	// Check ownership unless admin
	if role != models.RoleAdmin {
		var domainUserID int64
		h.db.QueryRow("SELECT user_id FROM domains WHERE id = ?", id).Scan(&domainUserID)
		if domainUserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error:   "Permission denied",
			})
		}
	}

	_, err = h.db.Exec("DELETE FROM domains WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to delete domain",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Domain deleted successfully",
	})
}
