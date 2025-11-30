package api

import (
	"strconv"

	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ListPackages(c *fiber.Ctx) error {
	rows, err := h.db.Query(`
		SELECT id, name, disk_quota, bandwidth_quota, max_domains, max_databases, max_emails, max_ftp, created_at
		FROM packages ORDER BY name
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch packages",
		})
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.Name, &p.DiskQuota, &p.BandwidthQuota, &p.MaxDomains, &p.MaxDatabases, &p.MaxEmails, &p.MaxFTP, &p.CreatedAt); err != nil {
			continue
		}
		packages = append(packages, p)
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    packages,
	})
}

func (h *Handler) CreatePackage(c *fiber.Ctx) error {
	var pkg models.Package
	if err := c.BodyParser(&pkg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if pkg.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Package name is required",
		})
	}

	result, err := h.db.Exec(`
		INSERT INTO packages (name, disk_quota, bandwidth_quota, max_domains, max_databases, max_emails, max_ftp)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, pkg.Name, pkg.DiskQuota, pkg.BandwidthQuota, pkg.MaxDomains, pkg.MaxDatabases, pkg.MaxEmails, pkg.MaxFTP)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Package name already exists",
		})
	}

	id, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Package created successfully",
		Data:    map[string]int64{"id": id},
	})
}

func (h *Handler) UpdatePackage(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid package ID",
		})
	}

	var pkg models.Package
	if err := c.BodyParser(&pkg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	_, err = h.db.Exec(`
		UPDATE packages SET name = ?, disk_quota = ?, bandwidth_quota = ?, 
		max_domains = ?, max_databases = ?, max_emails = ?, max_ftp = ?
		WHERE id = ?
	`, pkg.Name, pkg.DiskQuota, pkg.BandwidthQuota, pkg.MaxDomains, pkg.MaxDatabases, pkg.MaxEmails, pkg.MaxFTP, id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to update package",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Package updated successfully",
	})
}

func (h *Handler) DeletePackage(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid package ID",
		})
	}

	// Check if package is in use
	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM user_packages WHERE package_id = ?", id).Scan(&count)
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Cannot delete package that is in use",
		})
	}

	_, err = h.db.Exec("DELETE FROM packages WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to delete package",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Package deleted successfully",
	})
}
