package migration

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/asergenalkan/serverpanel/internal/config"
	"github.com/asergenalkan/serverpanel/internal/services/account"
	"github.com/asergenalkan/serverpanel/internal/services/dns"
	"github.com/asergenalkan/serverpanel/internal/webserver"
)

var (
	ErrInvalidBackup     = errors.New("invalid cPanel backup format")
	ErrExtractionFailed  = errors.New("backup extraction failed")
	ErrUserAlreadyExists = errors.New("username already exists")
)

// CPanelBackupInfo contains parsed information from a cPanel backup
type CPanelBackupInfo struct {
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Domain        string             `json:"domain"`
	PHPVersion    string             `json:"php_version"`
	Plan          string             `json:"plan"`
	CreatedAt     string             `json:"created_at"`
	HomeDir       string             `json:"home_dir"`
	DiskQuota     int64              `json:"disk_quota"`
	EmailLimit    int                `json:"email_limit"`
	HasNodejs     bool               `json:"has_nodejs"`
	NodejsVersion string             `json:"nodejs_version"`
	NodejsApps    []NodejsAppInfo    `json:"nodejs_apps"`
	Databases     []DatabaseInfo     `json:"databases"`
	EmailAccounts []EmailAccountInfo `json:"email_accounts"`
	FTPAccounts   []FTPAccountInfo   `json:"ftp_accounts"`
	DNSRecords    []DNSRecordInfo    `json:"dns_records"`
	Subdomains    []SubdomainInfo    `json:"subdomains"`
	CronJobs      []CronJobInfo      `json:"cron_jobs"`
	SSLCerts      []SSLCertInfo      `json:"ssl_certs"`
	DKIMKey       string             `json:"dkim_key"`
	BackupSize    int64              `json:"backup_size"`
	ExtractedPath string             `json:"extracted_path"`
}

type NodejsAppInfo struct {
	Name       string            `json:"name"`
	Path       string            `json:"path"`
	EntryPoint string            `json:"entry_point"`
	Version    string            `json:"version"`
	EnvVars    map[string]string `json:"env_vars"`
}

type DatabaseInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	HasDump  bool   `json:"has_dump"`
	DumpPath string `json:"dump_path"`
}

type EmailAccountInfo struct {
	Email    string `json:"email"`
	Domain   string `json:"domain"`
	HasMails bool   `json:"has_mails"`
}

type FTPAccountInfo struct {
	Username string `json:"username"`
	HomeDir  string `json:"home_dir"`
	HasHash  bool   `json:"has_hash"`
	Hash     string `json:"hash"`
}

type SubdomainInfo struct {
	Name         string `json:"name"`
	DocumentRoot string `json:"document_root"`
}

type CronJobInfo struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}

type SSLCertInfo struct {
	Domain      string `json:"domain"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

type DNSRecordInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// ImportOptions defines what to import from the backup
type ImportOptions struct {
	ImportFiles       bool   `json:"import_files"`
	ImportDatabases   bool   `json:"import_databases"`
	ImportEmails      bool   `json:"import_emails"`
	ImportDNS         bool   `json:"import_dns"`
	ImportFTP         bool   `json:"import_ftp"`
	ImportNodejs      bool   `json:"import_nodejs"`
	ImportCron        bool   `json:"import_cron"`
	ImportSSL         bool   `json:"import_ssl"`
	PackageID         int64  `json:"package_id"`
	NewPassword       string `json:"new_password"`
	OverwriteExisting bool   `json:"overwrite_existing"`
}

// ImportResult contains the result of the import operation
type ImportResult struct {
	Success  bool     `json:"success"`
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Domain   string   `json:"domain"`
	Imported []string `json:"imported"`
	Warnings []string `json:"warnings"`
	Errors   []string `json:"errors"`
}

// DB interface for database operations
type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
}

type Service struct {
	db  DB
	cfg *config.Config
}

func NewService(db DB) *Service {
	return &Service{
		db:  db,
		cfg: config.Get(),
	}
}

// ExtractAndAnalyze extracts the cPanel backup and analyzes its contents
func (s *Service) ExtractAndAnalyze(backupPath string) (*CPanelBackupInfo, error) {
	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "cpanel-backup-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	log.Printf("📦 Extracting cPanel backup to: %s", tempDir)

	// Extract tar.gz
	if err := s.extractTarGz(backupPath, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	// Find the backup root directory (usually named like the backup file)
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	var backupRoot string
	for _, entry := range entries {
		if entry.IsDir() {
			backupRoot = filepath.Join(tempDir, entry.Name())
			break
		}
	}

	if backupRoot == "" {
		backupRoot = tempDir
	}

	log.Printf("📂 Backup root: %s", backupRoot)

	// Parse backup info
	info, err := s.parseBackupInfo(backupRoot)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	// Get backup file size
	if stat, err := os.Stat(backupPath); err == nil {
		info.BackupSize = stat.Size()
	}

	info.ExtractedPath = backupRoot

	return info, nil
}

// extractTarGz extracts a tar.gz file to the destination directory
func (s *Service) extractTarGz(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)

		// Prevent path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dst)) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			os.Symlink(header.Linkname, target)
		}
	}

	return nil
}

// parseBackupInfo parses the backup directory and extracts account information
func (s *Service) parseBackupInfo(backupRoot string) (*CPanelBackupInfo, error) {
	info := &CPanelBackupInfo{
		NodejsApps:    []NodejsAppInfo{},
		Databases:     []DatabaseInfo{},
		EmailAccounts: []EmailAccountInfo{},
		FTPAccounts:   []FTPAccountInfo{},
		DNSRecords:    []DNSRecordInfo{},
		Subdomains:    []SubdomainInfo{},
		CronJobs:      []CronJobInfo{},
		SSLCerts:      []SSLCertInfo{},
	}

	// Parse cp/<username> file for account info
	cpDir := filepath.Join(backupRoot, "cp")
	if entries, err := os.ReadDir(cpDir); err == nil && len(entries) > 0 {
		cpFile := filepath.Join(cpDir, entries[0].Name())
		info.Username = entries[0].Name()
		s.parseCPFile(cpFile, info)
	}

	// Parse userdata for domain and PHP info
	userdataDir := filepath.Join(backupRoot, "userdata")
	s.parseUserdata(userdataDir, info)

	// Check for Node.js apps
	s.parseNodejsApps(backupRoot, info)

	// Parse MySQL databases
	s.parseMySQLDatabases(backupRoot, info)

	// Parse email accounts
	s.parseEmailAccounts(backupRoot, info)

	// Parse FTP accounts
	s.parseFTPAccounts(backupRoot, info)

	// Parse DNS zone
	s.parseDNSZone(backupRoot, info)

	// Parse cron jobs
	s.parseCronJobs(backupRoot, info)

	// Parse DKIM keys
	s.parseDKIMKeys(backupRoot, info)

	// Parse subdomains
	s.parseSubdomains(backupRoot, info)

	return info, nil
}

// parseCPFile parses the cp/<username> file
func (s *Service) parseCPFile(path string, info *CPanelBackupInfo) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "CONTACTEMAIL":
			info.Email = value
		case "DNS":
			if info.Domain == "" {
				info.Domain = value
			}
		case "PLAN":
			info.Plan = value
		case "STARTDATE":
			if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.CreatedAt = fmt.Sprintf("%d", ts)
			}
		case "MAX_EMAIL_PER_HOUR":
			if limit, err := strconv.Atoi(value); err == nil {
				info.EmailLimit = limit
			}
		case "DISK_BLOCK_LIMIT":
			if quota, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.DiskQuota = quota
			}
		}
	}
}

// parseUserdata parses the userdata directory
func (s *Service) parseUserdata(userdataDir string, info *CPanelBackupInfo) {
	entries, err := os.ReadDir(userdataDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "cache.json" || entry.Name() == "main" || entry.Name() == "scope" {
			continue
		}

		// Skip SSL and php-fpm files
		if strings.Contains(entry.Name(), "_SSL") || strings.Contains(entry.Name(), "php-fpm") {
			continue
		}

		// This is likely the main domain file
		domainFile := filepath.Join(userdataDir, entry.Name())
		content, err := os.ReadFile(domainFile)
		if err != nil {
			continue
		}

		// Parse YAML-like key: value format
		domainData := parseSimpleYAML(string(content))

		if servername, ok := domainData["servername"]; ok {
			info.Domain = servername
		}

		if phpversion, ok := domainData["phpversion"]; ok {
			// Convert ea-php81 to 8.1
			info.PHPVersion = strings.TrimPrefix(phpversion, "ea-php")
			if len(info.PHPVersion) >= 2 {
				info.PHPVersion = info.PHPVersion[:1] + "." + info.PHPVersion[1:]
			}
		}

		if homedir, ok := domainData["homedir"]; ok {
			info.HomeDir = homedir
		}

		break // Only need the first domain file
	}
}

// parseNodejsApps checks for Node.js applications
func (s *Service) parseNodejsApps(backupRoot string, info *CPanelBackupInfo) {
	homeDir := filepath.Join(backupRoot, "homedir")
	nodevenvDir := filepath.Join(homeDir, "nodevenv")

	log.Printf("🔍 Node.js app aranıyor - homeDir: %s, nodevenvDir: %s", homeDir, nodevenvDir)

	// Track found apps to avoid duplicates
	foundApps := make(map[string]bool)

	// Method 1: Check nodevenv directory (cPanel Node.js apps)
	if _, err := os.Stat(nodevenvDir); err == nil {
		log.Printf("✅ nodevenv klasörü bulundu")
		entries, _ := os.ReadDir(nodevenvDir)
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			appName := entry.Name()
			log.Printf("📦 nodevenv'de app bulundu: %s", appName)
			// Try to find the actual app directory
			possiblePaths := []string{
				filepath.Join(homeDir, appName),
				filepath.Join(homeDir, "public_html", appName),
			}
			for _, appPath := range possiblePaths {
				log.Printf("🔎 App yolu kontrol ediliyor: %s", appPath)
				if app := s.parseNodejsApp(appPath, appName, nodevenvDir); app != nil {
					log.Printf("✅ Node.js app tespit edildi: %s (path: %s, entry: %s, version: %s)", app.Name, app.Path, app.EntryPoint, app.Version)
					info.NodejsApps = append(info.NodejsApps, *app)
					foundApps[appPath] = true
					info.HasNodejs = true
					break
				}
			}
		}
	} else {
		log.Printf("⚠️ nodevenv klasörü bulunamadı: %v", err)
	}

	// Method 2: Scan homedir for any package.json files (standalone Node.js apps)
	entries, err := os.ReadDir(homeDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip common non-app directories
		name := entry.Name()
		if name == "mail" || name == "etc" || name == "logs" || name == "ssl" ||
			name == "tmp" || name == ".trash" || name == ".cpanel" || name == "nodevenv" ||
			strings.HasPrefix(name, ".") {
			continue
		}

		appPath := filepath.Join(homeDir, name)
		if foundApps[appPath] {
			continue
		}

		if app := s.parseNodejsApp(appPath, name, nodevenvDir); app != nil {
			info.NodejsApps = append(info.NodejsApps, *app)
			info.HasNodejs = true
		}
	}

	// Set default Node.js version if found apps but no version detected
	if info.HasNodejs && info.NodejsVersion == "" {
		info.NodejsVersion = "18"
	}
}

// parseNodejsApp parses a single Node.js application directory
func (s *Service) parseNodejsApp(appPath, appName, nodevenvDir string) *NodejsAppInfo {
	packageJsonPath := filepath.Join(appPath, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return nil
	}

	app := &NodejsAppInfo{
		Name: appName,
		Path: appName,
	}

	// Parse package.json
	if content, err := os.ReadFile(packageJsonPath); err == nil {
		var pkgJson map[string]interface{}
		if json.Unmarshal(content, &pkgJson) == nil {
			// Get app name from package.json if available
			if pkgName, ok := pkgJson["name"].(string); ok && pkgName != "" {
				app.Name = pkgName
			}
			if main, ok := pkgJson["main"].(string); ok {
				app.EntryPoint = main
			}
			if scripts, ok := pkgJson["scripts"].(map[string]interface{}); ok {
				if start, ok := scripts["start"].(string); ok && app.EntryPoint == "" {
					// Try to extract entry point from start script
					if strings.Contains(start, "server.js") {
						app.EntryPoint = "server.js"
					} else if strings.Contains(start, "app.js") {
						app.EntryPoint = "app.js"
					} else if strings.Contains(start, "index.js") {
						app.EntryPoint = "index.js"
					} else if strings.Contains(start, "next") {
						app.EntryPoint = "server.js" // Next.js custom server
					}
				}
			}
		}
	}

	// Check Node.js version from nodevenv
	if nodevenvDir != "" {
		versionDir := filepath.Join(nodevenvDir, appName)
		if vEntries, err := os.ReadDir(versionDir); err == nil {
			for _, ve := range vEntries {
				if ve.IsDir() && ve.Name() != ".lock" {
					app.Version = ve.Name()
					break
				}
			}
		}
	}

	// Parse .env file
	envPath := filepath.Join(appPath, ".env")
	if content, err := os.ReadFile(envPath); err == nil {
		app.EnvVars = make(map[string]string)
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
				app.EnvVars[key] = value
			}
		}
	}

	return app
}

// parseMySQLDatabases parses MySQL database dumps
func (s *Service) parseMySQLDatabases(backupRoot string, info *CPanelBackupInfo) {
	mysqlDir := filepath.Join(backupRoot, "mysql")
	if _, err := os.Stat(mysqlDir); os.IsNotExist(err) {
		return
	}

	// Look for .sql files ONLY in mysql directory (not recursive)
	entries, err := os.ReadDir(mysqlDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Skip system databases and non-sql files
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		if name == "mysql.sql" || name == "information_schema.sql" || name == "performance_schema.sql" {
			continue
		}

		fi, err := entry.Info()
		if err != nil {
			continue
		}

		dbName := strings.TrimSuffix(name, ".sql")
		dumpPath := filepath.Join(mysqlDir, name)

		// Verify it's a real database dump (check file size and content)
		if fi.Size() < 100 {
			// Too small, probably not a real dump - check content
			content, _ := os.ReadFile(dumpPath)
			contentStr := string(content)
			// Skip if it's just GRANT statements or empty
			if !strings.Contains(contentStr, "CREATE TABLE") && !strings.Contains(contentStr, "INSERT INTO") {
				continue
			}
		}

		info.Databases = append(info.Databases, DatabaseInfo{
			Name:     dbName,
			Size:     fi.Size(),
			HasDump:  true,
			DumpPath: dumpPath,
		})
	}
}

// parseEmailAccounts parses email accounts from the backup
func (s *Service) parseEmailAccounts(backupRoot string, info *CPanelBackupInfo) {
	mailDir := filepath.Join(backupRoot, "homedir", "mail")
	if _, err := os.Stat(mailDir); os.IsNotExist(err) {
		return
	}

	entries, err := os.ReadDir(mailDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		domain := entry.Name()
		domainMailDir := filepath.Join(mailDir, domain)

		// Check for email accounts in this domain
		accounts, err := os.ReadDir(domainMailDir)
		if err != nil {
			continue
		}

		for _, acc := range accounts {
			if !acc.IsDir() || strings.HasPrefix(acc.Name(), ".") {
				continue
			}

			info.EmailAccounts = append(info.EmailAccounts, EmailAccountInfo{
				Email:    acc.Name() + "@" + domain,
				Domain:   domain,
				HasMails: true,
			})
		}
	}
}

// parseFTPAccounts parses FTP accounts from proftpdpasswd
func (s *Service) parseFTPAccounts(backupRoot string, info *CPanelBackupInfo) {
	ftpFile := filepath.Join(backupRoot, "proftpdpasswd")
	content, err := os.ReadFile(ftpFile)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 6 {
			continue
		}

		info.FTPAccounts = append(info.FTPAccounts, FTPAccountInfo{
			Username: parts[0],
			HomeDir:  parts[5],
			HasHash:  len(parts[1]) > 0,
			Hash:     parts[1],
		})
	}
}

// parseDNSZone parses DNS zone file
func (s *Service) parseDNSZone(backupRoot string, info *CPanelBackupInfo) {
	dnsDir := filepath.Join(backupRoot, "dnszones")
	entries, err := os.ReadDir(dnsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".db") {
			continue
		}

		zonePath := filepath.Join(dnsDir, entry.Name())
		content, err := os.ReadFile(zonePath)
		if err != nil {
			continue
		}

		// Parse DNS records
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, ";") || line == "" || strings.HasPrefix(line, "$") {
				continue
			}

			// Simple DNS record parsing
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			record := DNSRecordInfo{
				Name: fields[0],
			}

			// Find TTL and type
			for i := 1; i < len(fields)-1; i++ {
				if ttl, err := strconv.Atoi(fields[i]); err == nil {
					record.TTL = ttl
					continue
				}
				if fields[i] == "IN" {
					continue
				}
				// This should be the record type
				record.Type = fields[i]
				record.Content = strings.Join(fields[i+1:], " ")
				break
			}

			if record.Type != "" && record.Type != "SOA" && record.Type != "NS" {
				info.DNSRecords = append(info.DNSRecords, record)
			}
		}
	}
}

// parseCronJobs parses cron jobs
func (s *Service) parseCronJobs(backupRoot string, info *CPanelBackupInfo) {
	cronDir := filepath.Join(backupRoot, "cron")
	entries, err := os.ReadDir(cronDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		cronPath := filepath.Join(cronDir, entry.Name())
		content, err := os.ReadFile(cronPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			// Parse cron line (5 time fields + command)
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}

			info.CronJobs = append(info.CronJobs, CronJobInfo{
				Schedule: strings.Join(fields[:5], " "),
				Command:  strings.Join(fields[5:], " "),
			})
		}
	}
}

// parseDKIMKeys parses DKIM keys
func (s *Service) parseDKIMKeys(backupRoot string, info *CPanelBackupInfo) {
	dkimDir := filepath.Join(backupRoot, "domainkeys", "private")
	entries, err := os.ReadDir(dkimDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		keyPath := filepath.Join(dkimDir, entry.Name())
		content, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}

		info.DKIMKey = string(content)
		break
	}
}

// parseSubdomains parses subdomains
func (s *Service) parseSubdomains(backupRoot string, info *CPanelBackupInfo) {
	sdsFile := filepath.Join(backupRoot, "sds2")
	content, err := os.ReadFile(sdsFile)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) >= 2 {
			info.Subdomains = append(info.Subdomains, SubdomainInfo{
				Name:         parts[0],
				DocumentRoot: parts[1],
			})
		}
	}
}

// Import imports the cPanel backup with the given options
func (s *Service) Import(info *CPanelBackupInfo, options ImportOptions) (*ImportResult, error) {
	result := &ImportResult{
		Username: info.Username,
		Domain:   info.Domain,
		Imported: []string{},
		Warnings: []string{},
		Errors:   []string{},
	}

	// Check if username already exists
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", info.Username).Scan(&count); err == nil && count > 0 {
		if !options.OverwriteExisting {
			return nil, ErrUserAlreadyExists
		}
		result.Warnings = append(result.Warnings, "Mevcut kullanıcı üzerine yazılacak")
	}

	// Create account using account service
	accountSvc := account.NewService(s.db)

	password := options.NewPassword
	if password == "" {
		password = generateRandomPassword()
	}

	acc, err := accountSvc.CreateAccount(account.CreateAccountRequest{
		Username:  info.Username,
		Email:     info.Email,
		Password:  password,
		Domain:    info.Domain,
		PackageID: options.PackageID,
	})

	if err != nil {
		// If user exists or email constraint failed, try to continue
		errStr := err.Error()
		if strings.Contains(errStr, "exists") || strings.Contains(errStr, "UNIQUE constraint") {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Kullanıcı/email zaten mevcut: %v", err))
			// Try to get existing user
			s.db.QueryRow("SELECT id FROM users WHERE username = ?", info.Username).Scan(&result.UserID)
		} else {
			result.Errors = append(result.Errors, fmt.Sprintf("Hesap oluşturulamadı: %v", err))
			return result, err
		}
	}

	if acc != nil {
		result.UserID = acc.ID
		result.Imported = append(result.Imported, "✅ Kullanıcı hesabı oluşturuldu")
	}

	// Get user ID if not set
	if result.UserID == 0 {
		s.db.QueryRow("SELECT id FROM users WHERE username = ?", info.Username).Scan(&result.UserID)
	}

	homeDir := filepath.Join(s.cfg.HomeBaseDir, info.Username)

	// Import files
	if options.ImportFiles {
		if err := s.importFiles(info, homeDir); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Dosyalar import edilemedi: %v", err))
		} else {
			result.Imported = append(result.Imported, "✅ Dosyalar import edildi")
		}
	}

	// Import DNS records
	if options.ImportDNS && len(info.DNSRecords) > 0 {
		if err := s.importDNSRecords(info); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("DNS kayıtları import edilemedi: %v", err))
		} else {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ %d DNS kaydı import edildi", len(info.DNSRecords)))
		}
	}

	// Import FTP accounts
	if options.ImportFTP && len(info.FTPAccounts) > 0 {
		imported := s.importFTPAccounts(info, result.UserID)
		if imported > 0 {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ %d FTP hesabı import edildi", imported))
		}
	}

	// Import Node.js apps
	if options.ImportNodejs && info.HasNodejs && len(info.NodejsApps) > 0 {
		imported := s.importNodejsApps(info, result.UserID, homeDir, result)
		if imported > 0 {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ %d Node.js uygulaması import edildi", imported))
		}
	}

	// Import databases
	if options.ImportDatabases && len(info.Databases) > 0 {
		imported := s.importDatabases(info, result.UserID)
		if imported > 0 {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ %d veritabanı import edildi", imported))
		}
	}

	// Import cron jobs
	if options.ImportCron && len(info.CronJobs) > 0 {
		imported := s.importCronJobs(info, result.UserID)
		if imported > 0 {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ %d cron job import edildi", imported))
		}
	}

	// Set PHP version if different from default
	if info.PHPVersion != "" && info.PHPVersion != s.cfg.PHPVersion {
		if err := s.setPHPVersion(info, result.UserID); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("PHP versiyonu ayarlanamadı: %v", err))
		} else {
			result.Imported = append(result.Imported, fmt.Sprintf("✅ PHP %s ayarlandı", info.PHPVersion))
		}
	}

	result.Success = len(result.Errors) == 0

	// Cleanup extracted files
	if info.ExtractedPath != "" {
		go func() {
			os.RemoveAll(filepath.Dir(info.ExtractedPath))
		}()
	}

	return result, nil
}

// importFiles copies files from backup to home directory
func (s *Service) importFiles(info *CPanelBackupInfo, homeDir string) error {
	srcHomeDir := filepath.Join(info.ExtractedPath, "homedir")

	// Copy public_html
	srcPublicHTML := filepath.Join(srcHomeDir, "public_html")
	dstPublicHTML := filepath.Join(homeDir, "public_html")

	if _, err := os.Stat(srcPublicHTML); err == nil {
		if config.IsDevelopment() {
			log.Printf("🔧 [SIMÜLASYON] cp -r %s/* %s/", srcPublicHTML, dstPublicHTML)
		} else {
			cmd := exec.Command("cp", "-r", srcPublicHTML+"/.", dstPublicHTML+"/")
			if err := cmd.Run(); err != nil {
				return err
			}
			// Fix ownership
			exec.Command("chown", "-R", info.Username+":"+info.Username, dstPublicHTML).Run()
		}
	}

	// Copy Node.js app directories
	for _, app := range info.NodejsApps {
		srcAppDir := filepath.Join(srcHomeDir, app.Path)
		dstAppDir := filepath.Join(homeDir, app.Path)

		if _, err := os.Stat(srcAppDir); err == nil {
			if config.IsDevelopment() {
				log.Printf("🔧 [SIMÜLASYON] cp -r %s %s", srcAppDir, dstAppDir)
			} else {
				os.MkdirAll(filepath.Dir(dstAppDir), 0755)
				cmd := exec.Command("cp", "-r", srcAppDir, dstAppDir)
				cmd.Run()
				exec.Command("chown", "-R", info.Username+":"+info.Username, dstAppDir).Run()
			}
		}
	}

	return nil
}

// importDNSRecords imports DNS records
func (s *Service) importDNSRecords(info *CPanelBackupInfo) error {
	dnsManager := dns.NewManager(s.cfg.SimulateMode, s.cfg.SimulateBasePath)

	for _, record := range info.DNSRecords {
		// Skip some cPanel specific records
		if strings.Contains(record.Name, "cpanel") || strings.Contains(record.Name, "whm") ||
			strings.Contains(record.Name, "webdisk") || strings.Contains(record.Name, "cpcontacts") ||
			strings.Contains(record.Name, "cpcalendars") {
			continue
		}

		name := strings.TrimSuffix(record.Name, info.Domain+".")
		name = strings.TrimSuffix(name, ".")

		if err := dnsManager.AddRecord(info.Domain, record.Type, name, strings.Trim(record.Content, `"`), record.TTL); err != nil {
			log.Printf("Warning: Failed to add DNS record %s: %v", record.Name, err)
		}
	}

	return nil
}

// importFTPAccounts imports FTP accounts
func (s *Service) importFTPAccounts(info *CPanelBackupInfo, userID int64) int {
	imported := 0
	homeDir := filepath.Join(s.cfg.HomeBaseDir, info.Username)

	for _, ftp := range info.FTPAccounts {
		// Skip the main account FTP
		if ftp.Username == info.Username {
			continue
		}

		ftpHomeDir := ftp.HomeDir
		if strings.HasPrefix(ftpHomeDir, "/home/"+info.Username) {
			ftpHomeDir = strings.Replace(ftpHomeDir, "/home/"+info.Username, homeDir, 1)
		}

		_, err := s.db.Exec(`
			INSERT INTO ftp_accounts (user_id, username, home_dir, active, created_at)
			VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
		`, userID, ftp.Username, ftpHomeDir)

		if err == nil {
			imported++
		}
	}

	return imported
}

// importNodejsApps imports Node.js applications
func (s *Service) importNodejsApps(info *CPanelBackupInfo, userID int64, homeDir string, result *ImportResult) int {
	imported := 0

	for _, app := range info.NodejsApps {
		appRoot := filepath.Join(homeDir, app.Path)
		startupFile := app.EntryPoint
		if startupFile == "" {
			startupFile = "server.js"
		}

		nodeVersion := app.Version
		if nodeVersion == "" {
			nodeVersion = "18"
		}

		// Get domain ID and domain name for app_url (primary domain)
		var domainID int64
		var domainName string
		s.db.QueryRow("SELECT id, name FROM domains WHERE user_id = ? AND domain_type = 'primary'", userID).Scan(&domainID, &domainName)

		// Assign a port (start from 3000)
		var maxPort int
		s.db.QueryRow("SELECT COALESCE(MAX(port), 2999) FROM nodejs_apps WHERE user_id = ?", userID).Scan(&maxPort)
		port := maxPort + 1

		// Build app URL (public URL without port - reverse proxy handles routing)
		appURL := ""
		if domainName != "" {
			appURL = fmt.Sprintf("https://%s", domainName)
		}

		// Convert env vars to JSON string
		envJSON := "{}"
		if len(app.EnvVars) > 0 {
			if envBytes, err := json.Marshal(app.EnvVars); err == nil {
				envJSON = string(envBytes)
			}
		}

		_, err := s.db.Exec(`
			INSERT INTO nodejs_apps (user_id, domain_id, name, app_root, startup_file, node_version, port, app_url, mode, environment, auto_restart, status, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'production', ?, 1, 'stopped', CURRENT_TIMESTAMP)
		`, userID, domainID, app.Name, appRoot, startupFile, nodeVersion, port, appURL, envJSON)

		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Node.js app eklenemedi (%s): %v", app.Name, err))
		} else {
			imported++
		}
	}

	return imported
}

// importDatabases imports MySQL databases
func (s *Service) importDatabases(info *CPanelBackupInfo, userID int64) int {
	imported := 0

	for _, db := range info.Databases {
		if !db.HasDump {
			continue
		}

		// Create database entry
		_, err := s.db.Exec(`
			INSERT INTO databases (user_id, name, type, created_at)
			VALUES (?, ?, 'mysql', CURRENT_TIMESTAMP)
		`, userID, db.Name)

		if err == nil {
			// Import the dump if not in development
			if !config.IsDevelopment() {
				mysqlPass := os.Getenv("MYSQL_ROOT_PASSWORD")
				if mysqlPass != "" {
					// Create the database
					exec.Command("mysql", "-uroot", "-p"+mysqlPass, "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", db.Name)).Run()
					// Import the dump
					cmd := exec.Command("mysql", "-uroot", "-p"+mysqlPass, db.Name)
					f, _ := os.Open(db.DumpPath)
					if f != nil {
						cmd.Stdin = f
						cmd.Run()
						f.Close()
					}
				}
			}
			imported++
		}
	}

	return imported
}

// importCronJobs imports cron jobs
func (s *Service) importCronJobs(info *CPanelBackupInfo, userID int64) int {
	imported := 0

	for _, cron := range info.CronJobs {
		// Parse schedule
		parts := strings.Fields(cron.Schedule)
		if len(parts) != 5 {
			continue
		}

		_, err := s.db.Exec(`
			INSERT INTO cron_jobs (user_id, name, command, minute, hour, day, month, weekday, active, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP)
		`, userID, "Imported Cron", cron.Command, parts[0], parts[1], parts[2], parts[3], parts[4])

		if err == nil {
			imported++
		}
	}

	return imported
}

// setPHPVersion sets PHP version for the domain
func (s *Service) setPHPVersion(info *CPanelBackupInfo, userID int64) error {
	// Check if PHP version is installed
	phpVersion := strings.Replace(info.PHPVersion, ".", "", 1)
	phpFpmSock := fmt.Sprintf("/run/php/php%s-fpm.sock", info.PHPVersion)

	if _, err := os.Stat(phpFpmSock); os.IsNotExist(err) {
		return fmt.Errorf("PHP %s kurulu değil", info.PHPVersion)
	}

	// Update domain PHP version
	_, err := s.db.Exec(`
		UPDATE domains SET php_version = ? WHERE user_id = ?
	`, phpVersion, userID)

	if err != nil {
		return err
	}

	// Update Apache vhost
	homeDir := filepath.Join(s.cfg.HomeBaseDir, info.Username)
	documentRoot := filepath.Join(homeDir, "public_html")

	driver := webserver.NewDriver(webserver.DriverApache, s.cfg.SimulateMode, s.cfg.SimulateBasePath)
	vhostConfig := webserver.VhostConfig{
		Domain:       info.Domain,
		Aliases:      []string{fmt.Sprintf("www.%s", info.Domain)},
		Username:     info.Username,
		DocumentRoot: documentRoot,
		HomeDir:      homeDir,
		PHPVersion:   info.PHPVersion,
	}

	return driver.CreateVhost(vhostConfig)
}

// Cleanup removes the extracted backup files
func (s *Service) Cleanup(extractedPath string) {
	if extractedPath != "" {
		os.RemoveAll(filepath.Dir(extractedPath))
	}
}

// generateRandomPassword generates a random password
func generateRandomPassword() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, 16)
	for i := range password {
		password[i] = chars[i%len(chars)]
	}
	return string(password)
}

// ValidatePHPVersion checks if the required PHP version is available
func (s *Service) ValidatePHPVersion(version string) (bool, string) {
	if version == "" {
		return true, s.cfg.PHPVersion
	}

	phpFpmSock := fmt.Sprintf("/run/php/php%s-fpm.sock", version)
	if _, err := os.Stat(phpFpmSock); os.IsNotExist(err) {
		return false, version
	}

	return true, version
}

// GetInstalledPHPVersions returns list of installed PHP versions
func (s *Service) GetInstalledPHPVersions() []string {
	versions := []string{}

	pattern := regexp.MustCompile(`php(\d+\.\d+)-fpm\.sock`)
	entries, err := os.ReadDir("/run/php")
	if err != nil {
		return versions
	}

	for _, entry := range entries {
		if matches := pattern.FindStringSubmatch(entry.Name()); len(matches) > 1 {
			versions = append(versions, matches[1])
		}
	}

	return versions
}

// parseSimpleYAML parses a simple YAML-like key: value format
func parseSimpleYAML(content string) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Parse key: value
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			result[key] = value
		}
	}

	return result
}
