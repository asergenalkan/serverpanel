package api

import (
	"strconv"

	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ListDatabases(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	role := c.Locals("role").(string)

	var query string
	var args []interface{}

	if role == models.RoleAdmin {
		query = `SELECT id, user_id, name, type, size, created_at FROM databases ORDER BY name`
	} else {
		query = `SELECT id, user_id, name, type, size, created_at FROM databases WHERE user_id = ? ORDER BY name`
		args = append(args, userID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch databases",
		})
	}
	defer rows.Close()

	var databases []models.Database
	for rows.Next() {
		var db models.Database
		if err := rows.Scan(&db.ID, &db.UserID, &db.Name, &db.Type, &db.Size, &db.CreatedAt); err != nil {
			continue
		}
		databases = append(databases, db)
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    databases,
	})
}

func (h *Handler) CreateDatabase(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
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
			Error:   "Database name is required",
		})
	}

	if req.Type == "" {
		req.Type = "mysql"
	}

	result, err := h.db.Exec(`
		INSERT INTO databases (user_id, name, type)
		VALUES (?, ?, ?)
	`, userID, req.Name, req.Type)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Database name already exists",
		})
	}

	id, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Database created successfully",
		Data:    map[string]int64{"id": id},
	})
}

func (h *Handler) DeleteDatabase(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid database ID",
		})
	}

	userID := c.Locals("user_id").(int64)
	role := c.Locals("role").(string)

	// Check ownership unless admin
	if role != models.RoleAdmin {
		var dbUserID int64
		h.db.QueryRow("SELECT user_id FROM databases WHERE id = ?", id).Scan(&dbUserID)
		if dbUserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error:   "Permission denied",
			})
		}
	}

	_, err = h.db.Exec("DELETE FROM databases WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to delete database",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Database deleted successfully",
	})
}
