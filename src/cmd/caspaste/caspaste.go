// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/casjay-forks/caspaste/src/internal/apiv1"
	"github.com/casjay-forks/caspaste/src/internal/cli"
	"github.com/casjay-forks/caspaste/src/internal/config"
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/raw"
	"github.com/casjay-forks/caspaste/src/internal/service"
	"github.com/casjay-forks/caspaste/src/internal/storage"
	"github.com/casjay-forks/caspaste/src/internal/web"
)

var Version = "unknown"

// getVersion reads version from release.txt or returns default
func getVersion() string {
	// If Version was set at build time (via -ldflags), use it
	if Version != "unknown" {
		return Version
	}

	// Try to read from release.txt
	data, err := os.ReadFile("release.txt")
	if err == nil {
		version := strings.TrimSpace(string(data))
		if version != "" {
			return version
		}
	}

	// Default version
	return "1.0.0"
}

func readFile(path string) (string, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file
	fileByte, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Return result
	return string(fileByte), nil
}

func exitOnError(e error) {
	fmt.Fprintln(os.Stderr, "error:", e.Error())
	os.Exit(1)
}

// ensureDirectories creates all necessary directories if they don't exist
func ensureDirectories(dataDir, configDir string) error {
	// Create data directory structure
	if dataDir != "" {
		dirs := []string{
			dataDir,
			dataDir + "/db",
			dataDir + "/backups",
		}

		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}

	// Create config directory
	if configDir != "" {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	return nil
}

// handleServiceCommand processes --service flag commands
func handleServiceCommand(command, address, dbSource, dataDir, configDir string) {
	// Get executable path
	executable, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	// Build service config
	svcConfig := service.ServiceConfig{
		Name:        "caspaste",
		DisplayName: "CasPaste Pastebin Service",
		Description: "Self-hosted pastebin service",
		Executable:  executable,
		Args:        buildServiceArgs(address, dbSource, dataDir, configDir),
		WorkingDir:  dataDir,
		User:        "caspaste",
	}

	mgr := service.New(svcConfig)

	switch command {
	case "start":
		err = mgr.Start()
	case "stop":
		err = mgr.Stop()
	case "restart":
		err = mgr.Restart()
	case "reload":
		err = mgr.Reload()
	case "--install", "install":
		err = mgr.Install()
	case "--uninstall", "uninstall":
		err = mgr.Uninstall()
	case "--disable", "disable":
		err = mgr.Disable()
	case "--help", "help":
		printServiceHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown service command: %s\n", command)
		printServiceHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Service operation failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// buildServiceArgs creates the argument list for service configuration
func buildServiceArgs(address, dbSource, dataDir, configDir string) []string {
	args := []string{}

	if address != "" && address != ":80" {
		args = append(args, "--address", address)
	}
	if dbSource != "" {
		args = append(args, "--db-source", dbSource)
	}
	if dataDir != "" {
		args = append(args, "--data", dataDir)
	}
	if configDir != "" {
		args = append(args, "--config", configDir)
	}

	return args
}

// printServiceHelp shows service command help
func printServiceHelp() {
	fmt.Println("CasPaste Service Management")
	fmt.Println("===========================")
	fmt.Println()
	fmt.Println("Usage: caspaste --service COMMAND")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  start        - Start the service")
	fmt.Println("  stop         - Stop the service")
	fmt.Println("  restart      - Restart the service")
	fmt.Println("  reload       - Reload service configuration")
	fmt.Println("  --install    - Install service for automatic startup")
	fmt.Println("  --uninstall  - Remove service")
	fmt.Println("  --disable    - Disable service from starting at boot")
	fmt.Println("  --help       - Show this help")
	fmt.Println()
}

// handleMaintenanceCommand processes --maintenance flag commands
func handleMaintenanceCommand(command, dbDriver, dbSource, dataDir, configDir string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "Maintenance command required\n")
		printMaintenanceHelp()
		os.Exit(1)
	}

	action := parts[0]
	var arg string
	if len(parts) > 1 {
		arg = parts[1]
	}

	switch action {
	case "backup":
		err := performBackup(dbDriver, dbSource, dataDir, configDir, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	case "restore":
		err := performRestore(dbDriver, dbSource, dataDir, configDir, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Restore failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	case "mode":
		if arg == "" {
			fmt.Fprintf(os.Stderr, "Mode argument required: enabled or disabled\n")
			os.Exit(1)
		}
		err := setMaintenanceMode(dataDir, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set maintenance mode: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "Unknown maintenance command: %s\n", action)
		printMaintenanceHelp()
		os.Exit(1)
	}
}

// printMaintenanceHelp shows maintenance command help
func printMaintenanceHelp() {
	fmt.Println("CasPaste Maintenance Mode")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("Usage: caspaste --maintenance COMMAND [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  backup [filename]         - Full disaster recovery backup (default: backup-YYYYMMDD-HHMMSS.tar.gz)")
	fmt.Println("  restore [filename]        - Restore from backup (default: latest backup)")
	fmt.Println("  mode {enabled|disabled}   - Enable or disable maintenance mode")
	fmt.Println()
	fmt.Println("Backup includes:")
	fmt.Println("  - Config directory (caspaste.yml and all config files)")
	fmt.Println("  - Data directory (db/caspaste.db and all data)")
	fmt.Println("  - External SQLite database (if located outside data_dir/db/)")
	fmt.Println()
	fmt.Println("Note: When using PostgreSQL/MariaDB, db/caspaste.db is a synchronized cache")
	fmt.Println("      that's included in backups for instant disaster recovery.")
	fmt.Println()
}

// checkAndMigrateDatabase checks if database driver/source changed and auto-migrates if needed
func checkAndMigrateDatabase(dataDir, configDir, newDriver, newSource string) error {
	stateFile := dataDir + "/.db-state"

	// Read previous database state if exists
	oldStateData, err := os.ReadFile(stateFile)
	var oldDriver, oldSource string
	if err == nil {
		parts := strings.SplitN(string(oldStateData), "\n", 2)
		if len(parts) >= 1 {
			oldDriver = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			oldSource = strings.TrimSpace(parts[1])
		}
	}

	// Normalize driver names for comparison
	normalizedNew := normalizeDriverName(newDriver)
	normalizedOld := normalizeDriverName(oldDriver)

	// If driver changed, perform automatic migration
	if oldDriver != "" && oldSource != "" && (normalizedOld != normalizedNew || oldSource != newSource) {
		fmt.Println()
		fmt.Println("⚠️  Database configuration change detected!")
		fmt.Printf("Old: %s (%s)\n", oldDriver, oldSource)
		fmt.Printf("New: %s (%s)\n", newDriver, newSource)
		fmt.Println()
		fmt.Println("Starting automatic database migration...")
		fmt.Println("This may take a few minutes depending on database size.")
		fmt.Println()

		// Create backup before migration
		backupFilename := "pre-migration-" + time.Now().Format("20060102-150405") + ".tar.gz"
		fmt.Printf("Creating safety backup: %s\n", dataDir+"/backups/"+backupFilename)
		performBackup(oldDriver, oldSource, dataDir, configDir, backupFilename)

		// Perform migration
		err := storage.MigrateDatabase(oldDriver, oldSource, newDriver, newSource)
		if err != nil {
			fmt.Println()
			fmt.Println("❌ Migration failed!")
			fmt.Printf("Error: %v\n", err)
			fmt.Println()
			fmt.Println("Your old database is still intact. To restore:")
			fmt.Printf("  caspaste --maintenance \"restore %s\" --data %s\n", backupFilename, dataDir)
			return fmt.Errorf("automatic migration failed")
		}

		fmt.Println()
		fmt.Println("✅ Migration completed successfully!")
		fmt.Println()
	}

	// Save current database state for next startup
	stateData := newDriver + "\n" + newSource
	err = os.WriteFile(stateFile, []byte(stateData), 0644)
	if err != nil {
		fmt.Printf("Warning: failed to save database state: %v\n", err)
	}

	return nil
}

// normalizeDriverName normalizes driver names for comparison
func normalizeDriverName(driver string) string {
	driver = strings.ToLower(driver)
	if driver == "mariadb" {
		return "mysql"
	}
	return driver
}

// performBackup creates a full disaster recovery backup
func performBackup(dbDriver, dbSource, dataDir, configDir, filename string) error {
	if dataDir == "" {
		dataDir = "."
	}

	// Generate filename if not provided
	if filename == "" {
		filename = fmt.Sprintf("backup-%s.tar.gz", time.Now().Format("20060102-150405"))
	}

	// Ensure backup directory exists
	backupDir := dataDir + "/backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupPath := backupDir + "/" + filename

	fmt.Println("Creating disaster recovery backup...")
	fmt.Println("Backing up:")
	fmt.Printf("  - Config: %s\n", configDir)
	fmt.Printf("  - Data: %s\n", dataDir)

	// Check if database is outside data_dir/db
	expectedDbPath := dataDir + "/db/"
	dbIsExternal := false
	if !strings.HasPrefix(dbSource, expectedDbPath) && (dbDriver == "sqlite3" || dbDriver == "sqlite") {
		dbIsExternal = true
		fmt.Printf("  - Database: %s (external)\n", dbSource)
	}

	fmt.Printf("Destination: %s\n", backupPath)
	fmt.Println()

	// Create temporary directory for staging backup
	tempDir := dataDir + "/.backup-temp"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Copy data directory
	cmd := exec.Command("cp", "-r", dataDir, tempDir+"/data")
	cmd.Run()

	// Copy config directory if exists
	if configDir != "" {
		if _, err := os.Stat(configDir); err == nil {
			cmd = exec.Command("cp", "-r", configDir, tempDir+"/config")
			cmd.Run()
		}
	}

	// Copy external database if needed
	if dbIsExternal {
		os.MkdirAll(tempDir+"/external-db", 0755)
		cmd = exec.Command("cp", dbSource, tempDir+"/external-db/caspaste.db")
		cmd.Run()
	}

	// Create tar.gz archive
	cmd = exec.Command("tar", "-czf", backupPath,
		"--exclude=backups",
		"--exclude=.backup-temp",
		"--exclude=*.tmp",
		"--exclude=*.lock",
		"-C", tempDir,
		".")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("backup failed: %w\nOutput: %s", err, string(output))
	}

	// Get backup file size
	info, err := os.Stat(backupPath)
	if err == nil {
		fmt.Printf("Backup created: %s (%.2f MB)\n", backupPath, float64(info.Size())/1024/1024)
	} else {
		fmt.Printf("Backup created: %s\n", backupPath)
	}

	return nil
}

// performRestore performs full disaster recovery restore from backup archive
func performRestore(dbDriver, dbSource, dataDir, configDir, filename string) error {
	if dataDir == "" {
		dataDir = "."
	}

	backupDir := dataDir + "/backups"

	// If no filename, find latest backup
	if filename == "" {
		entries, err := os.ReadDir(backupDir)
		if err != nil {
			return fmt.Errorf("failed to read backup directory: %w", err)
		}

		var latestFile string
		var latestTime int64

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tar.gz") {
				info, err := entry.Info()
				if err != nil {
					continue
				}
				if info.ModTime().Unix() > latestTime {
					latestTime = info.ModTime().Unix()
					latestFile = entry.Name()
				}
			}
		}

		if latestFile == "" {
			return fmt.Errorf("no backup files found in %s", backupDir)
		}

		filename = latestFile
		fmt.Printf("Using latest backup: %s\n", filename)
	}

	backupPath := backupDir + "/" + filename

	// Check backup exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}

	// Create safety backup of current state
	fmt.Println("Creating safety backup of current state...")
	performBackup(dbDriver, dbSource, dataDir, configDir, "pre-restore-"+time.Now().Format("20060102-150405")+".tar.gz")

	// Create temporary extraction directory
	tempDir := dataDir + "/.restore-temp"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Extract backup archive to temp directory
	fmt.Printf("Restoring from: %s\n", backupPath)
	fmt.Println("Extracting backup archive...")

	cmd := exec.Command("tar", "-xzf", backupPath, "-C", tempDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restore failed: %w\nOutput: %s", err, string(output))
	}

	// Restore data directory
	if _, err := os.Stat(tempDir + "/data"); err == nil {
		fmt.Println("Restoring data directory...")
		cmd = exec.Command("cp", "-r", tempDir+"/data/.", dataDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to restore data directory: %w", err)
		}
	}

	// Restore config directory
	if configDir != "" {
		if _, err := os.Stat(tempDir + "/config"); err == nil {
			fmt.Println("Restoring config directory...")
			cmd = exec.Command("cp", "-r", tempDir+"/config/.", configDir)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to restore config directory: %w", err)
			}
		}
	}

	// Restore external database if exists
	if _, err := os.Stat(tempDir + "/external-db/caspaste.db"); err == nil {
		fmt.Println("Restoring external database...")
		if dbDriver == "sqlite3" || dbDriver == "sqlite" {
			cmd = exec.Command("cp", tempDir+"/external-db/caspaste.db", dbSource)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to restore external database: %w", err)
			}
		}
	}

	fmt.Println()
	fmt.Println("Disaster recovery restore completed successfully")
	fmt.Println("Restored:")
	fmt.Printf("  - Data: %s\n", dataDir)
	if configDir != "" {
		fmt.Printf("  - Config: %s\n", configDir)
	}
	return nil
}

// setMaintenanceMode enables or disables maintenance mode
func setMaintenanceMode(dataDir, mode string) error {
	maintenanceFile := dataDir
	if maintenanceFile == "" {
		maintenanceFile = "."
	}
	maintenanceFile = maintenanceFile + "/.maintenance"

	switch mode {
	case "enabled", "enable", "on":
		err := os.WriteFile(maintenanceFile, []byte("maintenance mode enabled"), 0644)
		if err != nil {
			return fmt.Errorf("failed to enable maintenance mode: %w", err)
		}
		fmt.Println("Maintenance mode: ENABLED")
		fmt.Printf("Maintenance file created: %s\n", maintenanceFile)
		return nil

	case "disabled", "disable", "off":
		if err := os.Remove(maintenanceFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to disable maintenance mode: %w", err)
		}
		fmt.Println("Maintenance mode: DISABLED")
		return nil

	default:
		return fmt.Errorf("invalid mode: %s (use 'enabled' or 'disabled')", mode)
	}
}

// checkStatus performs health check on database and returns exit code
// Exit codes: 0 = healthy, 1 = unhealthy, 2 = error
func checkStatus(dbDriver, dbSource string, address string) {
	fmt.Println("CasPaste Health Check")
	fmt.Println("=====================")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Listen Address: %s\n", address)
	fmt.Printf("Database Driver: %s\n", dbDriver)
	fmt.Println()

	exitCode := 0
	healthy := true

	// Check database connectivity
	fmt.Print("Checking database connection... ")
	db, err := storage.NewPool(dbDriver, dbSource, 1, 0, "")
	if err != nil {
		fmt.Printf("FAILED\n  Error: %v\n", err)
		healthy = false
		exitCode = 1
	} else {
		// Try to ping the database
		err = db.Close()
		if err != nil {
			fmt.Printf("DEGRADED\n  Warning: %v\n", err)
			exitCode = 2
		} else {
			fmt.Println("OK")
		}
	}

	// Check if we can initialize database schema
	if healthy {
		fmt.Print("Checking database schema... ")
		err = storage.InitDB(dbDriver, dbSource)
		if err != nil {
			fmt.Printf("FAILED\n  Error: %v\n", err)
			healthy = false
			exitCode = 1
		} else {
			fmt.Println("OK")
		}
	}

	fmt.Println()
	if healthy && exitCode == 0 {
		fmt.Println("Status: HEALTHY")
		os.Exit(0)
	} else if exitCode == 2 {
		fmt.Println("Status: DEGRADED")
		os.Exit(2)
	} else {
		fmt.Println("Status: UNHEALTHY")
		os.Exit(1)
	}
}

func main() {
	var err error

	// Get version (from build-time, release.txt, or default)
	Version = getVersion()

	// Read environment variables and CLI flags
	c := cli.New(Version)

	flagAddress := c.AddStringVar("address", ":80", "HTTP server ADDRESS:PORT (use FQDN for reverse proxy setups).", &cli.FlagOptions{
		PreHook: func(s string) (string, error) {
			if s == "" {
				return s, nil
			}

			// If the address doesn't contain a colon, it's missing the port
			if !strings.Contains(s, ":") {
				// Check if it looks like a FQDN (contains a dot)
				if strings.Contains(s, ".") {
					// FQDN without port: bind to all interfaces on port 80
					// The actual public URL will be constructed using reverse proxy headers
					return ":80", nil
				}
				// IP address or hostname without port: append :80
				return s + ":80", nil
			}

			// Check if it's a FQDN with port (e.g., "example.com:8080")
			parts := strings.Split(s, ":")
			if len(parts) == 2 && strings.Contains(parts[0], ".") {
				// FQDN with port: bind to all interfaces on the specified port
				// The actual public URL will be constructed using reverse proxy headers
				return ":" + parts[1], nil
			}

			return s, nil
		},
	})

	// Special commands (don't require full setup)
	flagStatus := c.AddBoolVar("status", "Check server health and database connectivity. Exit codes: 0=healthy, 1=unhealthy, 2=error")
	flagService := c.AddStringVar("service", "", "Service management: start, stop, restart, reload, --install, --uninstall, --disable, --help", nil)
	flagMaintenance := c.AddStringVar("maintenance", "", "Maintenance mode: backup [filename], restore [filename], mode {enabled|disabled}", nil)

	// New modern flags (also support old flag names for backward compatibility)
	flagPort := c.AddStringVar("port", "", "Port to listen on (alternative to specifying in --address). Examples: 80, 8080, 443.", nil)
	flagDataDir := c.AddStringVar("data", "", "Data directory for storing database and other files.", nil)
	flagConfigDir := c.AddStringVar("config", "", "Configuration directory for loading config files.", nil)

	flagDbDriver := c.AddStringVar("db-driver", "sqlite3", "Database driver: sqlite3, postgres, mysql, mariadb", nil)
	flagDbSource := c.AddStringVar("db-source", "", "DB source (auto-set when using --data).", nil)
	flagDbMaxOpenConns := c.AddIntVar("db-max-open-conns", 25, "Maximum number of connections to the database.", nil)
	flagDbMaxIdleConns := c.AddIntVar("db-max-idle-conns", 5, "Maximum number of idle connections to the database.", nil)
	flagDbCleanupPeriod := c.AddDurationVar("db-cleanup-period", "1m", "Interval at which the DB is cleared of expired but not yet deleted pastes.", nil)

	flagRobotsDisallow := c.AddBoolVar("robots-disallow", "Prohibits search engine crawlers from indexing site using robots.txt file.")

	flagTitleMaxLen := c.AddIntVar("title-max-length", 100, "Maximum length of the paste title. If 0 disable title, if -1 disable length limit.", nil)
	flagBodyMaxLen := c.AddIntVar("body-max-length", 52428800, "Maximum length of the paste body in bytes. Default 50MB. If -1 disable length limit. Can't be -1.", nil)
	flagMaxLifetime := c.AddDurationVar("max-paste-lifetime", "unlimited", "Maximum lifetime of the paste. Examples: 10m, 1h 30m, 12h, 1w, 30d, 365d.", &cli.FlagOptions{
		PreHook: func(s string) (string, error) {
			if s == "never" || s == "unlimited" {
				return "", nil
			}

			return s, nil
		},
	})

	flagGetPastesPer5Min := c.AddUintVar("get-pastes-per-5min", 50, "Maximum number of pastes that can be VIEWED in 5 minutes from one IP. If 0 disable rate-limit.", nil)
	flagGetPastesPer15Min := c.AddUintVar("get-pastes-per-15min", 100, "Maximum number of pastes that can be VIEWED in 15 minutes from one IP. If 0 disable rate-limit.", nil)
	flagGetPastesPer1Hour := c.AddUintVar("get-pastes-per-1hour", 500, "Maximum number of pastes that can be VIEWED in 1 hour from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer5Min := c.AddUintVar("new-pastes-per-5min", 15, "Maximum number of pastes that can be CREATED in 5 minutes from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer15Min := c.AddUintVar("new-pastes-per-15min", 30, "Maximum number of pastes that can be CREATED in 15 minutes from one IP. If 0 disable rate-limit.", nil)
	flagNewPastesPer1Hour := c.AddUintVar("new-pastes-per-1hour", 40, "Maximum number of pastes that can be CREATED in 1 hour from one IP. If 0 disable rate-limit.", nil)

	flagServerAbout := c.AddStringVar("server-about", "", "Path to the TXT file that contains the server description.", nil)
	flagServerRules := c.AddStringVar("server-rules", "", "Path to the TXT file that contains the server rules.", nil)
	flagServerTerms := c.AddStringVar("server-terms", "", "Path to the TXT file that contains the server terms of use.", nil)

	flagAdminName := c.AddStringVar("admin-name", "", "Name of the administrator of this server.", nil)
	flagAdminMail := c.AddStringVar("admin-mail", "", "Email of the administrator of this server.", nil)

	flagUiDefaultLifetime := c.AddStringVar("ui-default-lifetime", "never", "Lifetime of paste will be set by default in WEB interface. Examples: 10min, 1h, 1d, 2w, 6mon, 1y, never.", nil)
	flagUiDefaultTheme := c.AddStringVar("ui-default-theme", "dracula", "Sets the default theme for the WEB interface. Examples: dracula, nord, github-light.", nil)
	flagUiThemesDir := c.AddStringVar("ui-themes-dir", "", "Loads external WEB interface themes from directory.", nil)

	flagCasPasswdFile := c.AddStringVar("caspasswd-file", "", "File in CasPasswd format. If set, authorization will be required to create pastes.", nil)

	flagTrustReverseProxy := c.AddBoolVar("trust-reverse-proxy", "Trust X-Forwarded-* headers for client IP detection. Only enable when behind a trusted reverse proxy (nginx, caddy, etc.). WARNING: Enabling without a proxy allows IP spoofing.")

	c.Parse()

	// Try to load config file from config directory or current directory
	var yamlCfg *config.YAMLConfig
	configPaths := []string{}
	if *flagConfigDir != "" {
		configPaths = append(configPaths, *flagConfigDir+"/caspaste.yml", *flagConfigDir+"/caspaste.yaml")
	}
	configPaths = append(configPaths, "caspaste.yml", "caspaste.yaml", "/etc/caspaste/caspaste.yml", "/etc/caspaste/caspaste.yaml")

	for _, path := range configPaths {
		cfg, err := config.LoadYAMLConfig(path)
		if err == nil {
			yamlCfg = cfg
			fmt.Printf("Loaded config from: %s\n", path)
			break
		}
	}

	// If no config file found and --config specified, create default config
	if yamlCfg == nil && *flagConfigDir != "" {
		defaultConfigPath := *flagConfigDir + "/caspaste.yml"
		if err := config.GenerateDefaultYAMLConfig(defaultConfigPath); err == nil {
			fmt.Printf("Created default config file: %s\n", defaultConfigPath)
			// Try to load the newly created config
			if cfg, err := config.LoadYAMLConfig(defaultConfigPath); err == nil {
				yamlCfg = cfg
			}
		}
	}

	// Merge config file with flags (flags take precedence)
	if yamlCfg != nil {
		if *flagPort == "" && yamlCfg.Server.Port != 0 {
			*flagPort = strconv.Itoa(yamlCfg.Server.Port)
		}
		if *flagAddress == ":80" && yamlCfg.Server.Address != "" {
			*flagAddress = yamlCfg.Server.Address
		}
		if *flagDbSource == "" && yamlCfg.Database.Source != "" {
			*flagDbSource = yamlCfg.Database.Source
		}
		if *flagDbDriver == "sqlite3" && yamlCfg.Database.Driver != "" {
			*flagDbDriver = yamlCfg.Database.Driver
		}
		if *flagAdminName == "" && yamlCfg.Server.AdminName != "" {
			*flagAdminName = yamlCfg.Server.AdminName
		}
		if *flagAdminMail == "" && yamlCfg.Server.AdminEmail != "" {
			*flagAdminMail = yamlCfg.Server.AdminEmail
		}
		if *flagUiDefaultTheme == "dracula" && yamlCfg.UI.DefaultTheme != "" {
			*flagUiDefaultTheme = yamlCfg.UI.DefaultTheme
		}
		if *flagCasPasswdFile == "" && yamlCfg.Security.PasswordFile != "" {
			*flagCasPasswdFile = yamlCfg.Security.PasswordFile
		}
	}

	// Process --port flag (overrides port in --address)
	if *flagPort != "" {
		// Extract host from address (if any)
		addr := *flagAddress
		if strings.Contains(addr, ":") {
			// Remove existing port
			parts := strings.Split(addr, ":")
			addr = parts[0]
		}
		// Append new port
		if !strings.HasPrefix(*flagPort, ":") {
			*flagAddress = addr + ":" + *flagPort
		} else {
			*flagAddress = addr + *flagPort
		}
	}

	// Process --data directory
	if *flagDataDir != "" {
		// Update db-source if it's relative or default
		if *flagDbSource == "" || !strings.Contains(*flagDbSource, "/") {
			*flagDbSource = *flagDataDir + "/db/caspaste.db"
		}
	}

	// Ensure all directories exist
	if err := ensureDirectories(*flagDataDir, *flagConfigDir); err != nil {
		exitOnError(err)
	}

	// Process --config directory (for future config file support)
	if *flagConfigDir != "" {
		// Reserved for future use - could load config files from this directory
		// For now, just note it for potential caspasswd-file path resolution
		if *flagCasPasswdFile != "" && !strings.HasPrefix(*flagCasPasswdFile, "/") {
			*flagCasPasswdFile = *flagConfigDir + "/" + *flagCasPasswdFile
		}
	}

	// Validate that either --data or --db-source is provided
	if *flagDbSource == "" {
		exitOnError(errors.New("either --data or --db-source must be provided"))
	}

	// Handle --status command (exits after checking)
	if *flagStatus {
		checkStatus(*flagDbDriver, *flagDbSource, *flagAddress)
		return // checkStatus calls os.Exit, but return for safety
	}

	// Handle --service command (exits after operation)
	if *flagService != "" {
		handleServiceCommand(*flagService, *flagAddress, *flagDbSource, *flagDataDir, *flagConfigDir)
		return
	}

	// Handle --maintenance command (exits after operation)
	if *flagMaintenance != "" {
		handleMaintenanceCommand(*flagMaintenance, *flagDbDriver, *flagDbSource, *flagDataDir, *flagConfigDir)
		return
	}

	// Auto-detect and perform database migration if driver changed
	if *flagDataDir != "" {
		err := checkAndMigrateDatabase(*flagDataDir, *flagConfigDir, *flagDbDriver, *flagDbSource)
		if err != nil {
			exitOnError(err)
		}
	}

	// -body-max-length flag
	if *flagBodyMaxLen == 0 {
		exitOnError(errors.New("maximum body length cannot be 0"))
	}

	// -max-paste-lifetime
	maxLifeTime := int64(-1)

	if *flagMaxLifetime != 0 && *flagMaxLifetime < 600 {
		exitOnError(errors.New("maximum paste lifetime flag cannot have a value less than 10 minutes"))
		maxLifeTime = int64(*flagMaxLifetime / time.Second)
	}

	// Load server about
	serverAbout := ""
	if *flagServerAbout != "" {
		serverAbout, err = readFile(*flagServerAbout)
		if err != nil {
			exitOnError(err)
		}
	}

	// Load server rules
	serverRules := ""
	if *flagServerRules != "" {
		serverRules, err = readFile(*flagServerRules)
		if err != nil {
			exitOnError(err)
		}
	}

	// Load server "terms of use"
	serverTermsOfUse := ""
	if *flagServerTerms != "" {
		if serverRules == "" {
			exitOnError(errors.New("in order to set the Terms of Use you must also specify the Server Rules"))
		}

		serverTermsOfUse, err = readFile(*flagServerTerms)
		if err != nil {
			exitOnError(err)
		}
	}

	// Settings
	log := logger.New("2006/01/02 15:04:05")

	db, err := storage.NewPool(*flagDbDriver, *flagDbSource, *flagDbMaxOpenConns, *flagDbMaxIdleConns, *flagDataDir)
	if err != nil {
		exitOnError(err)
	}

	cfg := config.Config{
		Log:               log,
		RateLimitGet:      netshare.NewRateLimitSystem(*flagGetPastesPer5Min, *flagGetPastesPer15Min, *flagGetPastesPer1Hour),
		RateLimitNew:      netshare.NewRateLimitSystem(*flagNewPastesPer5Min, *flagNewPastesPer15Min, *flagNewPastesPer1Hour),
		Version:           Version,
		TitleMaxLen:       *flagTitleMaxLen,
		BodyMaxLen:        *flagBodyMaxLen,
		MaxLifeTime:       maxLifeTime,
		ServerAbout:       serverAbout,
		ServerRules:       serverRules,
		ServerTermsOfUse:  serverTermsOfUse,
		AdminName:         *flagAdminName,
		AdminMail:         *flagAdminMail,
		RobotsDisallow:    *flagRobotsDisallow,
		TrustReverseProxy: *flagTrustReverseProxy,
		UiDefaultLifetime: *flagUiDefaultLifetime,
		UiDefaultTheme:    *flagUiDefaultTheme,
		UiThemesDir:       *flagUiThemesDir,
		CasPasswdFile:     *flagCasPasswdFile,
	}

	apiv1Data := apiv1.Load(db, cfg)

	rawData := raw.Load(db, cfg)

	// Init data base
	err = storage.InitDB(*flagDbDriver, *flagDbSource)
	if err != nil {
		exitOnError(err)
	}

	// Load pages
	webData, err := web.Load(db, cfg)
	if err != nil {
		exitOnError(err)
	}

	// Handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		webData.Handler(rw, req)
	})
	mux.HandleFunc("/raw/", func(rw http.ResponseWriter, req *http.Request) {
		rawData.Hand(rw, req)
	})
	mux.HandleFunc("/api/", func(rw http.ResponseWriter, req *http.Request) {
		apiv1Data.Hand(rw, req)
	})

	// Wrap with maintenance mode middleware
	dataDirectory := *flagDataDir
	if dataDirectory == "" {
		dataDirectory = "."
	}
	handler := web.MaintenanceMiddleware(dataDirectory, mux)

	// Run background job
	go func(cleanJobPeriod time.Duration) {
		for {
			// Delete expired pastes
			count, err := db.PasteDeleteExpired()
			if err != nil {
				log.Error(errors.New("Delete expired: " + err.Error()))
			}

			log.Info("Delete " + strconv.FormatInt(count, 10) + " expired pastes")

			// Wait
			time.Sleep(cleanJobPeriod)
		}
	}(*flagDbCleanupPeriod)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         *flagAddress,
		Handler:      handler, // Custom mux with maintenance middleware
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup signal handling for graceful shutdown
	// Works on Windows, macOS, BSD, and Linux
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Run HTTP server on " + *flagAddress)
		serverErrors <- srv.ListenAndServe()
	}()

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			exitOnError(err)
		}

	case sig := <-sigChan:
		log.Info(fmt.Sprintf("Received signal %v, shutting down gracefully...", sig))

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Error(fmt.Errorf("server shutdown error: %w", err))
			// Force close if graceful shutdown fails
			srv.Close()
		}

		log.Info("Server stopped")
	}
}
