package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate [source] [destination]",
		Short: "Migrate spirit state to new machine/server",
		Long: `Move your agent's spirit to a new home.

Supports migration between:
- Local directories
- Git repositories  
- Different backends (GitHub → S3, etc.)

Example:
  spirit migrate ~/old-spirit ~/new-spirit
  spirit migrate github:TheOrionAI/orion-state s3://my-bucket/orion`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrateSpirit(args[0], args[1])
		},
	}
}

func migrateSpirit(source, dest string) error {
	fmt.Printf("🌌 Migrating SPIRIT from '%s' to '%s'\n\n", source, dest)

	// Detect source type and parse
	sourceType, sourcePath := parseLocation(source)
	destType, destPath := parseLocation(dest)

	fmt.Printf("Source: %s (%s)\n", sourceType, sourcePath)
	fmt.Printf("Destination: %s (%s)\n\n", destType, destPath)

	// Export from source
	fmt.Println("📦 Exporting state...")
	exportData, err := exportFrom(sourceType, sourcePath)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	// Validate export
	if err := validateExport(exportData); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Import to destination
	fmt.Println("📥 Importing to new location...")
	if err := importTo(destType, destPath, exportData); err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	// Update local config to point to new location
	if destType == "github" || destType == "gitlab" {
		if err := updatePrimaryBackend(dest); err != nil {
			return fmt.Errorf("backend update failed: %w", err)
		}
	}

	fmt.Println("\n✅ Migration complete!")
	fmt.Printf("Your agent's spirit is now at: %s\n", dest)
	fmt.Println("\nTo verify:")
	fmt.Printf("  spirit status\n")

	return nil
}

func parseLocation(loc string) (string, string) {
	// Parse location strings like:
	// github:TheOrionAI/orion-state
	// s3://my-bucket/orion
	// /local/path
	// current (use ~/.spirit)

	if loc == "current" || loc == "." {
		return "local", ConfigDir
	}

	if len(loc) > 7 && loc[:7] == "github:" {
		return "github", loc[7:]
	}

	if len(loc) > 7 && loc[:7] == "gitlab:" {
		return "gitlab", loc[7:]
	}

	if len(loc) > 5 && loc[:5] == "s3://" {
		return "s3", loc[5:]
	}

	// Assume local path
	return "local", loc
}

type ExportPackage struct {
	Version   string                 `json:"version"`
	Identity  Identity               `json:"identity"`
	Soul      Soul                   `json:"soul"`
	Backends  map[string]Backend    `json:"backends"`
	Memory    []string               `json:"memory,omitempty"`
	Projects  []string               `json:"projects,omitempty"`
	ExportedAt string               `json:"exported_at"`
}

func exportFrom(sourceType, sourcePath string) (*ExportPackage, error) {
	var pkg ExportPackage

	switch sourceType {
	case "local":
		// Read from local filesystem
		configPath := filepath.Join(sourcePath, "spirit.json")
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read config: %w", err)
		}

		var config Config
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("cannot parse config: %w", err)
		}

		pkg.Identity = config.Identity
		pkg.Soul = config.Soul
		pkg.Backends = config.Backends

		// TODO: Read memory and projects

	case "github", "gitlab":
		// Clone and read
		// Implementation depends on git client
		return nil, fmt.Errorf("github export not yet implemented")

	case "s3":
		// Download from S3
		return nil, fmt.Errorf("s3 export not yet implemented")
	}

	// Capture export timestamp
	pkg.Version = "1.0.0"

	return &pkg, nil
}

func validateExport(pkg *ExportPackage) error {
	if pkg.Identity.Name == "" {
		return fmt.Errorf("identity missing name")
	}
	return nil
}

func importTo(destType, destPath string, pkg *ExportPackage) error {
	switch destType {
	case "local":
		// Ensure directory exists
		if err := os.MkdirAll(destPath, 0755); err != nil {
			return fmt.Errorf("cannot create directory: %w", err)
		}

		// Create subdirectories
		for _, dir := range []string{"memory", "projects", "context"} {
			if err := os.MkdirAll(filepath.Join(destPath, dir), 0755); err != nil {
				return fmt.Errorf("cannot create %s: %w", dir, err)
			}
		}

		// Write config
		config := Config{
			Version:   pkg.Version,
			Identity:  pkg.Identity,
			Soul:      pkg.Soul,
			Backends:  pkg.Backends,
		}

		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("cannot marshal config: %w", err)
		}

		configPath := filepath.Join(destPath, "spirit.json")
		if err := os.WriteFile(configPath, data, 0600); err != nil {
			return fmt.Errorf("cannot write config: %w", err)
		}

	case "github", "gitlab":
		// Create repo and push
		return fmt.Errorf("github import requires git - not yet implemented")

	case "s3":
		// Upload to S3
		return fmt.Errorf("s3 import not yet implemented")
	}

	return nil
}

func updatePrimaryBackend(newBackend string) error {
	// Update ~/.spirit/spirit.json primary backend
	configPath := filepath.Join(ConfigDir, "spirit.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	config.Backends["primary"] = Backend{
		Type: "github",
		Config: map[string]string{
			"repo": newBackend,
		},
	}

	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, newData, 0600)
}
