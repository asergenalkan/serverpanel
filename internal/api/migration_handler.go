package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/asergenalkan/serverpanel/internal/models"
	"github.com/asergenalkan/serverpanel/internal/services/migration"
	"github.com/gofiber/fiber/v2"
)

// UploadCPanelBackup handles cPanel backup upload and analysis
func (h *Handler) UploadCPanelBackup(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("backup")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Backup dosyası gerekli",
		})
	}

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".tar.gz") &&
		!strings.HasSuffix(strings.ToLower(file.Filename), ".tgz") {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Sadece .tar.gz veya .tgz dosyaları kabul edilir",
		})
	}

	// Create temp directory for upload
	tempDir, err := os.MkdirTemp("", "cpanel-upload-*")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Geçici dizin oluşturulamadı",
		})
	}

	// Save uploaded file
	tempPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveFile(file, tempPath); err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Dosya kaydedilemedi",
		})
	}

	// Analyze the backup
	svc := migration.NewService(h.db)
	info, err := svc.ExtractAndAnalyze(tempPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Backup analiz edilemedi: %v", err),
		})
	}

	// Store the temp path for later import
	info.ExtractedPath = info.ExtractedPath

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    info,
	})
}

// AnalyzeCPanelBackup analyzes an already uploaded backup
func (h *Handler) AnalyzeCPanelBackup(c *fiber.Ctx) error {
	var req struct {
		BackupPath string `json:"backup_path"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Geçersiz istek",
		})
	}

	if req.BackupPath == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Backup yolu gerekli",
		})
	}

	svc := migration.NewService(h.db)
	info, err := svc.ExtractAndAnalyze(req.BackupPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Backup analiz edilemedi: %v", err),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    info,
	})
}

// ImportCPanelBackup imports a cPanel backup with the given options
func (h *Handler) ImportCPanelBackup(c *fiber.Ctx) error {
	var req struct {
		BackupInfo    *migration.CPanelBackupInfo `json:"backup_info"`
		ImportOptions migration.ImportOptions     `json:"import_options"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Geçersiz istek formatı",
		})
	}

	if req.BackupInfo == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Backup bilgisi gerekli",
		})
	}

	if req.ImportOptions.PackageID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Paket seçimi gerekli",
		})
	}

	svc := migration.NewService(h.db)

	// Validate PHP version if needed
	if req.BackupInfo.PHPVersion != "" {
		available, version := svc.ValidatePHPVersion(req.BackupInfo.PHPVersion)
		if !available {
			// Check installed versions
			installed := svc.GetInstalledPHPVersions()
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error:   fmt.Sprintf("PHP %s kurulu değil. Kurulu versiyonlar: %v", version, installed),
			})
		}
	}

	// Perform import
	result, err := svc.Import(req.BackupInfo, req.ImportOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Import başarısız: %v", err),
			Data:    result,
		})
	}

	return c.JSON(models.APIResponse{
		Success: result.Success,
		Data:    result,
	})
}

// GetInstalledPHPVersions returns list of installed PHP versions for migration
func (h *Handler) GetMigrationPHPVersions(c *fiber.Ctx) error {
	svc := migration.NewService(h.db)
	versions := svc.GetInstalledPHPVersions()

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    versions,
	})
}

// CleanupMigration cleans up temporary files from a migration
func (h *Handler) CleanupMigration(c *fiber.Ctx) error {
	var req struct {
		ExtractedPath string `json:"extracted_path"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Geçersiz istek",
		})
	}

	if req.ExtractedPath != "" {
		svc := migration.NewService(h.db)
		svc.Cleanup(req.ExtractedPath)
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Temizlendi",
	})
}

// StreamUploadCPanelBackup handles large file uploads with streaming
func (h *Handler) StreamUploadCPanelBackup(c *fiber.Ctx) error {
	// Get file from multipart form
	file, err := c.FormFile("backup")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Backup dosyası gerekli",
		})
	}

	// Validate file extension
	filename := strings.ToLower(file.Filename)
	if !strings.HasSuffix(filename, ".tar.gz") && !strings.HasSuffix(filename, ".tgz") {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Sadece .tar.gz veya .tgz dosyaları kabul edilir",
		})
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cpanel-upload-*")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Geçici dizin oluşturulamadı",
		})
	}

	tempPath := filepath.Join(tempDir, file.Filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Dosya açılamadı",
		})
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(tempPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Hedef dosya oluşturulamadı",
		})
	}
	defer dst.Close()

	// Copy with progress tracking
	_, err = io.Copy(dst, src)
	if err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   "Dosya kopyalanamadı",
		})
	}

	// Analyze the backup
	svc := migration.NewService(h.db)
	info, err := svc.ExtractAndAnalyze(tempPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Backup analiz edilemedi: %v", err),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    info,
	})
}
