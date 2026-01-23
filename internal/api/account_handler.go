package api

import (
	"strconv"

	"github.com/asergenalkan/serverpanel/internal/auth"
	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/asergenalkan/serverpanel/internal/services/account"
	"github.com/gofiber/fiber/v2"
)

// ListAccounts returns all hosting accounts (Admin only)
func (h *Handler) ListAccounts(c *fiber.Ctx) error {
	svc := account.NewService(h.db)

	accounts, err := svc.ListAccounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Failed to fetch accounts",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    accounts,
	})
}

// CreateAccount creates a new hosting account (Admin only)
func (h *Handler) CreateAccount(c *fiber.Ctx) error {
	var req account.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Domain == "" || req.PackageID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "All fields are required: username, email, password, domain, package_id",
		})
	}

	svc := account.NewService(h.db)

	acc, err := svc.CreateAccount(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Data:    acc,
	})
}

// GetAccount returns a specific account
func (h *Handler) GetAccount(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid account ID",
		})
	}

	var acc struct {
		ID          int64  `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Role        string `json:"role"`
		Active      bool   `json:"active"`
		CreatedAt   string `json:"created_at"`
		Domain      string `json:"domain"`
		PackageName string `json:"package_name"`
	}

	err = h.db.QueryRow(`
		SELECT u.id, u.username, u.email, u.role, u.active, u.created_at,
			   COALESCE(d.name, '') as domain,
			   COALESCE(p.name, 'No Package') as package_name
		FROM users u
		LEFT JOIN domains d ON d.user_id = u.id
		LEFT JOIN user_packages up ON up.user_id = u.id
		LEFT JOIN packages p ON p.id = up.package_id
		WHERE u.id = ?
	`, id).Scan(&acc.ID, &acc.Username, &acc.Email, &acc.Role, &acc.Active,
		&acc.CreatedAt, &acc.Domain, &acc.PackageName)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   "Account not found",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    acc,
	})
}

// DeleteAccount deletes a hosting account
func (h *Handler) DeleteAccount(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid account ID",
		})
	}

	svc := account.NewService(h.db)

	if err := svc.DeleteAccount(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Account deleted successfully"},
	})
}

// SuspendAccount suspends an account
func (h *Handler) SuspendAccount(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid account ID",
		})
	}

	svc := account.NewService(h.db)

	if err := svc.SuspendAccount(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Account suspended"},
	})
}

// UnsuspendAccount unsuspends an account
func (h *Handler) UnsuspendAccount(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid account ID",
		})
	}

	svc := account.NewService(h.db)

	if err := svc.UnsuspendAccount(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Account unsuspended"},
	})
}

// ResetAccountPassword resets a hosting account's panel password (Admin only)
func (h *Handler) ResetAccountPassword(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid account ID",
		})
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Şifre en az 8 karakter olmalıdır",
		})
	}

	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND role = 'user'", id).Scan(&count)
	if count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   "Account not found",
		})
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Şifre oluşturma hatası",
		})
	}

	_, err = h.db.Exec("UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", hashed, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Şifre güncellenemedi",
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Şifre başarıyla güncellendi",
	})
}
