
// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/casjay-forks/caspaste/src/internal/apiv1"
	"github.com/casjay-forks/caspaste/src/internal/cli"
	"github.com/casjay-forks/caspaste/src/internal/config"
	"github.com/casjay-forks/caspaste/src/internal/logger"
	"github.com/casjay-forks/caspaste/src/internal/netshare"
	"github.com/casjay-forks/caspaste/src/internal/privilege"
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

// getDisplayAddress converts a listen address to a user-friendly display address
// Replaces 0.0.0.0, 127.0.0.1, localhost, etc. with valid FQDN, hostname, or IP
func getDisplayAddress(listenAddr string) string {
	host, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		// No port specified, use address as-is
		host = listenAddr
		port = "80"
	}

	// List of addresses to replace (localhost/loopback indicators)
	replaceableHosts := []string{"", "0.0.0.0", "127.0.0.1", "localhost", "::1", "::"}

	shouldReplace := false
	for _, replaceable := range replaceableHosts {
		if host == replaceable {
			shouldReplace = true
			break
		}
	}

	if shouldReplace {
		// Try to get hostname
		if hostname, err := os.Hostname(); err == nil && hostname != "" && hostname != "localhost" {
			host = hostname
		} else {
			// Try to get first non-loopback IP
			if addrs, err := net.InterfaceAddrs(); err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							// Prefer IPv4
							host = ipnet.IP.String()
							break
						}
					}
				}
			}
		}
	}

	// If still couldn't determine, use localhost
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}

	return net.JoinHostPort(host, port)
}

// isRunningAsRoot checks if the process is running with root/admin privileges
func isRunningAsRoot() bool {
	switch runtime.GOOS {
	case "windows":
		// On Windows, check if running as administrator
		// Simple heuristic: try to create a file in Windows system directory
		testPath := os.Getenv("WINDIR") + "\\Temp\\.caspaste-test"
		if f, err := os.Create(testPath); err == nil {
			f.Close()
			os.Remove(testPath)
			return true
		}
		return false
	default:
		// Unix-like systems: check if UID is 0
		return os.Geteuid() == 0
	}
}

// ensureDirectories creates all necessary directories if they don't exist
func ensureDirectories(dataDir, configDir, dbDir, backupDir, cacheDir, logsDir string) error {
	// Create data directory
	if dataDir != "" {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dataDir, err)
		}
	}

	// Create database directory if specified and different from dataDir
	if dbDir != "" && dbDir != dataDir {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dbDir, err)
		}
	}

	// Create backup directory if specified
	if backupDir != "" {
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", backupDir, err)
		}
	}

	// Create cache directory if specified
	if cacheDir != "" {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", cacheDir, err)
		}
	}

	// Create logs directory if specified
	if logsDir != "" {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", logsDir, err)
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
func handleMaintenanceCommand(command, dbDriver, dbSource, dataDir, configDir, backupDir string) {
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
		err := performBackup(dbDriver, dbSource, dataDir, configDir, backupDir, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	case "restore":
		err := performRestore(dbDriver, dbSource, dataDir, configDir, backupDir, arg)
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
func checkAndMigrateDatabase(dataDir, configDir, backupDir, newDriver, newSource string) error {
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
		fmt.Printf("Creating safety backup: %s\n", backupDir+"/"+backupFilename)
		performBackup(oldDriver, oldSource, dataDir, configDir, backupDir, backupFilename)

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

// normalizeDriverName normalizes driver names for comparison and usage
func normalizeDriverName(driver string) string {
	driver = strings.ToLower(driver)
	// MariaDB uses MySQL driver
	if driver == "mariadb" {
		return "mysql"
	}
	// sqlite3 (CGo driver) → sqlite (pure Go driver)
	// We use modernc.org/sqlite (pure Go) which registers as "sqlite"
	if driver == "sqlite3" {
		return "sqlite"
	}
	return driver
}

// performBackup creates a full disaster recovery backup
func performBackup(dbDriver, dbSource, dataDir, configDir, backupDir, filename string) error {
	if dataDir == "" {
		dataDir = "."
	}

	// Generate filename if not provided
	if filename == "" {
		filename = fmt.Sprintf("backup-%s.tar.gz", time.Now().Format("20060102-150405"))
	}

	// Ensure backup directory exists
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
func performRestore(dbDriver, dbSource, dataDir, configDir, backupDir, filename string) error {
	if dataDir == "" {
		dataDir = "."
	}

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
	performBackup(dbDriver, dbSource, dataDir, configDir, backupDir, "pre-restore-"+time.Now().Format("20060102-150405")+".tar.gz")

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

	// Directory flags
	flagPort := c.AddStringVar("port", "", "Port to listen on (alternative to specifying in --address). Examples: 80, 8080, 443.", nil)
	flagDataDir := c.AddStringVar("data", "", "Data directory. Examples: /var/lib/caspaste, ~/.local/share/caspaste", nil)
	flagConfigDir := c.AddStringVar("config", "", "Configuration directory. Examples: /etc/caspaste, ~/.config/caspaste", nil)
	flagCacheDir := c.AddStringVar("cache", "", "Cache directory. Examples: /var/cache/caspaste, ~/.cache/caspaste", nil)
	flagLogsDir := c.AddStringVar("logs", "", "Logs directory. Examples: /var/log/caspaste, ~/.local/log/caspaste", nil)

	c.Parse()

	// Create config directory first if specified (needed before generating config file)
	if *flagConfigDir != "" {
		if err := os.MkdirAll(*flagConfigDir, 0755); err != nil {
			exitOnError(fmt.Errorf("failed to create config directory: %w", err))
		}
	}

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

	// If no config file found, create default config
	if yamlCfg == nil {
		var defaultConfigPath string
		if *flagConfigDir != "" {
			defaultConfigPath = *flagConfigDir + "/caspaste.yml"
		} else {
			defaultConfigPath = "./caspaste.yml"
		}

		if err := config.GenerateDefaultYAMLConfig(defaultConfigPath); err != nil {
			exitOnError(fmt.Errorf("failed to create default config file: %w", err))
		}

		fmt.Printf("Created default config file: %s\n", defaultConfigPath)

		// Load the newly created config
		cfg, err := config.LoadYAMLConfig(defaultConfigPath)
		if err != nil {
			exitOnError(fmt.Errorf("failed to load generated config: %w", err))
		}
		yamlCfg = cfg
	}

	// Apply environment variable overrides to config
	// Priority: Config file < Environment variables < CLI flags
	config.ApplyEnvironmentOverrides(yamlCfg)

	// Merge CLI flags (highest priority - override both config file and env vars)
	if *flagPort != "" {
		yamlCfg.Server.Port, _ = strconv.Atoi(*flagPort)
	}
	if *flagAddress != ":80" {
		yamlCfg.Server.Address = *flagAddress
	} else if yamlCfg.Server.Address != "" {
		*flagAddress = yamlCfg.Server.Address
	}

	// Merge cache/logs directories from CLI (override config if specified)
	if *flagCacheDir != "" {
		yamlCfg.Directories.Cache = *flagCacheDir
	}
	if *flagLogsDir != "" {
		yamlCfg.Directories.Logs = *flagLogsDir
	}

	// Normalize database driver name (sqlite3 → sqlite, mariadb → mysql)
	yamlCfg.Database.Driver = normalizeDriverName(yamlCfg.Database.Driver)

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

	// Process --data directory and determine database directory from config
	var dbDir string
	dbSource := yamlCfg.Database.Source
	if dbSource == "" {
		exitOnError(errors.New("database.source must be specified in config file"))
	}

	// Only process file paths for SQLite databases
	// PostgreSQL/MySQL use connection strings (postgres://, mysql://, etc.)
	driver := yamlCfg.Database.Driver
	if driver == "sqlite" || driver == "sqlite3" {
		// If database source is relative, make it absolute based on data directory
		if !strings.HasPrefix(dbSource, "/") && *flagDataDir != "" {
			// Check for CASPASTE_DB_DIR or LENPASTE_DB_DIR environment variable
			// These specify the DATABASE DIRECTORY (not full path)
			dbDir = os.Getenv("CASPASTE_DB_DIR")
			if dbDir == "" {
				dbDir = os.Getenv("LENPASTE_DB_DIR") // Backward compatibility
			}
			if dbDir == "" {
				// Code default: {dataDir}/db
				// Docker overrides this via ENV CASPASTE_DB_DIR=/data/db/sqlite in Dockerfile
				dbDir = *flagDataDir + "/db"
			}
			yamlCfg.Database.Source = dbDir + "/caspaste.db"
			dbSource = yamlCfg.Database.Source
		}

		// Extract directory from database source path
		if strings.Contains(dbSource, "/") {
			lastSlash := strings.LastIndex(dbSource, "/")
			if lastSlash > 0 {
				dbDir = dbSource[:lastSlash]
			}
		}
	}

	// Determine backup directory
	backupDir := os.Getenv("CASPASTE_BACKUP_DIR")
	if backupDir == "" {
		backupDir = os.Getenv("LENPASTE_BACKUP_DIR") // Backward compatibility
	}
	if backupDir == "" && *flagDataDir != "" {
		// Set platform-specific defaults
		if *flagDataDir == "/data" {
			// Docker container
			backupDir = "/data/backups"
		} else {
			// Standalone binary - use platform-specific defaults
			// Prefer global system directories if running as root, fallback to user directories
			isRoot := isRunningAsRoot()

			switch runtime.GOOS {
			case "linux":
				if isRoot {
					// Root: Use /mnt/Backups/caspaste (global)
					backupDir = "/mnt/Backups/caspaste"
				} else {
					// User: Use ~/.local/share/caspaste/backups
					if home := os.Getenv("HOME"); home != "" {
						backupDir = home + "/.local/share/caspaste/backups"
					} else {
						backupDir = *flagDataDir + "/backups"
					}
				}

			case "darwin":
				if isRoot {
					// Root: Use /var/backups/caspaste (global)
					backupDir = "/var/backups/caspaste"
				} else {
					// User: Use ~/Library/Application Support/CasPaste/Backups
					if home := os.Getenv("HOME"); home != "" {
						backupDir = home + "/Library/Application Support/CasPaste/Backups"
					} else {
						backupDir = *flagDataDir + "/backups"
					}
				}

			case "windows":
				if isRoot {
					// Admin: Use C:\ProgramData\CasPaste\Backups (global)
					if programData := os.Getenv("ProgramData"); programData != "" {
						backupDir = programData + "\\CasPaste\\Backups"
					} else {
						backupDir = "C:\\ProgramData\\CasPaste\\Backups"
					}
				} else {
					// User: Use %APPDATA%\CasPaste\Backups
					if appdata := os.Getenv("APPDATA"); appdata != "" {
						backupDir = appdata + "\\CasPaste\\Backups"
					} else {
						backupDir = *flagDataDir + "/backups"
					}
				}

			case "freebsd", "openbsd":
				if isRoot {
					// Root: Use /var/backups/caspaste (global)
					backupDir = "/var/backups/caspaste"
				} else {
					// User: Use ~/.caspaste/backups
					if home := os.Getenv("HOME"); home != "" {
						backupDir = home + "/.caspaste/backups"
					} else {
						backupDir = *flagDataDir + "/backups"
					}
				}

			default:
				// Fallback
				backupDir = *flagDataDir + "/backups"
			}
		}
	}

	// Determine cache directory
	cacheDir := yamlCfg.Directories.Cache
	if cacheDir == "" {
		cacheDir = os.Getenv("CASPASTE_CACHE_DIR")
	}
	if cacheDir == "" {
		cacheDir = os.Getenv("LENPASTE_CACHE_DIR") // Backward compatibility
	}
	if cacheDir == "" && *flagDataDir != "" {
		isRoot := isRunningAsRoot()
		switch runtime.GOOS {
		case "linux":
			if isRoot {
				cacheDir = "/var/cache/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					cacheDir = home + "/.cache/caspaste"
				} else {
					cacheDir = *flagDataDir + "/cache"
				}
			}
		case "darwin":
			if isRoot {
				cacheDir = "/var/cache/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					cacheDir = home + "/Library/Caches/CasPaste"
				} else {
					cacheDir = *flagDataDir + "/cache"
				}
			}
		case "windows":
			if isRoot {
				cacheDir = "C:\\ProgramData\\CasPaste\\Cache"
			} else {
				if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
					cacheDir = localAppData + "\\CasPaste\\Cache"
				} else {
					cacheDir = *flagDataDir + "/cache"
				}
			}
		case "freebsd", "openbsd":
			if isRoot {
				cacheDir = "/var/cache/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					cacheDir = home + "/.cache/caspaste"
				} else {
					cacheDir = *flagDataDir + "/cache"
				}
			}
		default:
			cacheDir = *flagDataDir + "/cache"
		}
	}

	// Determine logs directory
	logsDir := yamlCfg.Directories.Logs
	if logsDir == "" {
		logsDir = os.Getenv("CASPASTE_LOGS_DIR")
	}
	if logsDir == "" {
		logsDir = os.Getenv("LENPASTE_LOGS_DIR") // Backward compatibility
	}
	if logsDir == "" && *flagDataDir != "" {
		isRoot := isRunningAsRoot()
		switch runtime.GOOS {
		case "linux":
			if isRoot {
				logsDir = "/var/log/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					logsDir = home + "/.local/log/caspaste"
				} else {
					logsDir = *flagDataDir + "/logs"
				}
			}
		case "darwin":
			if isRoot {
				logsDir = "/var/log/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					logsDir = home + "/Library/Logs/CasPaste"
				} else {
					logsDir = *flagDataDir + "/logs"
				}
			}
		case "windows":
			if isRoot {
				logsDir = "C:\\ProgramData\\CasPaste\\Logs"
			} else {
				if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
					logsDir = localAppData + "\\CasPaste\\Logs"
				} else {
					logsDir = *flagDataDir + "/logs"
				}
			}
		case "freebsd", "openbsd":
			if isRoot {
				logsDir = "/var/log/caspaste"
			} else {
				if home := os.Getenv("HOME"); home != "" {
					logsDir = home + "/.local/log/caspaste"
				} else {
					logsDir = *flagDataDir + "/logs"
				}
			}
		default:
			logsDir = *flagDataDir + "/logs"
		}
	}

	// Setup user (Linux/BSD/macOS only) - must be done before creating directories
	var uid, gid int
	if runtime.GOOS != "windows" {
		var err error
		uid, gid, err = privilege.EnsureUser()
		if err != nil {
			// User creation failed - might not be running as root or user already exists
			// This is OK - we'll create directories with current user
			uid = 0
			gid = 0
		}
	}

	// Ensure all directories exist
	if err := ensureDirectories(*flagDataDir, *flagConfigDir, dbDir, backupDir, cacheDir, logsDir); err != nil {
		exitOnError(err)
	}

	// Chown ALL directories if we're running as root and created a user
	// This must be done before privilege drop to ensure the user can access everything
	if os.Geteuid() == 0 && uid > 0 && gid > 0 {
		dirsToChown := []string{*flagDataDir, *flagConfigDir, dbDir, backupDir, cacheDir, logsDir}
		for _, dir := range dirsToChown {
			if dir != "" {
				if err := privilege.ChownPathRecursive(dir, uid, gid); err != nil {
					// Log but don't fail - directory might not exist or already has correct ownership
					fmt.Fprintf(os.Stderr, "Warning: failed to chown %s: %v\n", dir, err)
				}
			}
		}
	}

	// Handle --status command (exits after checking)
	if *flagStatus {
		checkStatus(yamlCfg.Database.Driver, yamlCfg.Database.Source, *flagAddress)
		return // checkStatus calls os.Exit, but return for safety
	}

	// Handle --service command (exits after operation)
	if *flagService != "" {
		handleServiceCommand(*flagService, *flagAddress, yamlCfg.Database.Source, *flagDataDir, *flagConfigDir)
		return
	}

	// Handle --maintenance command (exits after operation)
	if *flagMaintenance != "" {
		handleMaintenanceCommand(*flagMaintenance, yamlCfg.Database.Driver, yamlCfg.Database.Source, *flagDataDir, *flagConfigDir, backupDir)
		return
	}

	// Auto-detect and perform database migration if driver changed
	if *flagDataDir != "" {
		err := checkAndMigrateDatabase(*flagDataDir, *flagConfigDir, backupDir, yamlCfg.Database.Driver, yamlCfg.Database.Source)
		if err != nil {
			exitOnError(err)
		}
	}

	// Validate body max length from config
	if yamlCfg.Limits.BodyMaxLength == 0 {
		exitOnError(errors.New("limits.body_max_length cannot be 0 in config file"))
	}

	// Parse max paste lifetime from config
	maxLifeTime := int64(-1)
	if yamlCfg.Limits.MaxPasteLifetime != "" && yamlCfg.Limits.MaxPasteLifetime != "never" && yamlCfg.Limits.MaxPasteLifetime != "unlimited" {
		duration, err := cli.ParseDuration(yamlCfg.Limits.MaxPasteLifetime)
		if err != nil {
			exitOnError(fmt.Errorf("invalid limits.max_paste_lifetime in config: %w", err))
		}
		if duration < 600*time.Second {
			exitOnError(errors.New("limits.max_paste_lifetime cannot be less than 10 minutes"))
		}
		maxLifeTime = int64(duration / time.Second)
	}

	// Load server about from config
	serverAbout := ""
	if yamlCfg.Content.AboutFile != "" {
		serverAbout, err = readFile(yamlCfg.Content.AboutFile)
		if err != nil {
			exitOnError(fmt.Errorf("failed to read content.about_file: %w", err))
		}
	}

	// Load server rules from config
	serverRules := ""
	if yamlCfg.Content.RulesFile != "" {
		serverRules, err = readFile(yamlCfg.Content.RulesFile)
		if err != nil {
			exitOnError(fmt.Errorf("failed to read content.rules_file: %w", err))
		}
	}

	// Load server terms of use from config
	serverTermsOfUse := ""
	if yamlCfg.Content.TermsFile != "" {
		if serverRules == "" {
			exitOnError(errors.New("content.terms_file requires content.rules_file to also be set"))
		}
		serverTermsOfUse, err = readFile(yamlCfg.Content.TermsFile)
		if err != nil {
			exitOnError(fmt.Errorf("failed to read content.terms_file: %w", err))
		}
	}

	// Settings
	log := logger.New("2006/01/02 15:04:05")

	db, err := storage.NewPool(yamlCfg.Database.Driver, yamlCfg.Database.Source, yamlCfg.Database.MaxOpenConns, yamlCfg.Database.MaxIdleConns, *flagDataDir)
	if err != nil {
		exitOnError(err)
	}

	cfg := config.Config{
		Log:               log,
		RateLimitGet:      netshare.NewRateLimitSystem(yamlCfg.Limits.GetPastesPer5Min, yamlCfg.Limits.GetPastesPer15Min, yamlCfg.Limits.GetPastesPer1Hour),
		RateLimitNew:      netshare.NewRateLimitSystem(yamlCfg.Limits.NewPastesPer5Min, yamlCfg.Limits.NewPastesPer15Min, yamlCfg.Limits.NewPastesPer1Hour),
		Version:           Version,
		TitleMaxLen:       yamlCfg.Limits.TitleMaxLength,
		BodyMaxLen:        yamlCfg.Limits.BodyMaxLength,
		MaxLifeTime:       maxLifeTime,
		ServerAbout:       serverAbout,
		ServerRules:       serverRules,
		ServerTermsOfUse:  serverTermsOfUse,
		AdminName:         yamlCfg.Server.AdminName,
		AdminMail:         yamlCfg.Server.AdminEmail,
		RobotsDisallow:    yamlCfg.Server.RobotsDisallow,
		TrustReverseProxy: yamlCfg.Server.TrustReverseProxy,
		UiDefaultLifetime: yamlCfg.UI.DefaultLifetime,
		UiDefaultTheme:    yamlCfg.UI.DefaultTheme,
		UiThemesDir:       yamlCfg.UI.ThemesDir,
		CasPasswdFile:     yamlCfg.Security.PasswordFile,
	}

	apiv1Data := apiv1.Load(db, cfg)

	rawData := raw.Load(db, cfg)

	// Init data base
	err = storage.InitDB(yamlCfg.Database.Driver, yamlCfg.Database.Source)
	if err != nil {
		exitOnError(err)
	}

	// Chown directories AGAIN after database initialization to ensure DB file has correct ownership
	// The database file was just created, so it needs to be chowned before privilege drop
	if os.Geteuid() == 0 && uid > 0 && gid > 0 {
		dirsToChown := []string{*flagDataDir, *flagConfigDir, dbDir, backupDir, cacheDir, logsDir}
		for _, dir := range dirsToChown {
			if dir != "" {
				privilege.ChownPathRecursive(dir, uid, gid)
			}
		}
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
	// Parse cleanup period from config
	cleanupPeriod, err := cli.ParseDuration(yamlCfg.Database.CleanupPeriod)
	if err != nil {
		exitOnError(fmt.Errorf("invalid database.cleanup_period in config: %w", err))
	}

	// Apply middleware chain: CORS → Maintenance → App
	handler := web.CORSMiddleware(web.MaintenanceMiddleware(dataDirectory, mux))

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
	}(cleanupPeriod)

	// Create listener (must be done as root for ports < 1024 on Unix)
	listener, err := net.Listen("tcp", *flagAddress)
	if err != nil {
		exitOnError(fmt.Errorf("failed to bind to %s: %w", *flagAddress, err))
	}

	// Drop privileges after binding to port (uid/gid set earlier during directory creation)
	if runtime.GOOS != "windows" && uid > 0 && gid > 0 {
		if err := privilege.DropPrivileges(uid, gid); err != nil {
			log.Error(fmt.Errorf("failed to drop privileges: %w", err))
			// Continue anyway
		}
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Handler:      handler, // Custom mux with middleware
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup signal handling for graceful shutdown
	// Works on Windows, macOS, BSD, and Linux
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine using the already-created listener
	serverErrors := make(chan error, 1)
	go func() {
		displayAddr := getDisplayAddress(*flagAddress)
		log.Info("Run HTTP server on " + displayAddr)
		serverErrors <- srv.Serve(listener)
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
