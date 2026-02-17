package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type Identity struct {
	Name        string    `json:"name"`
	Emoji       string    `json:"emoji"`
	Email       string    `json:"email,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Soul struct {
	Vibe       string   `json:"vibe"`
	CoreTruths []string `json:"core_truths"`
	Boundaries []string `json:"boundaries"`
}

type Config struct {
	Version   string                 `json:"version"`
	Backends  map[string]Backend    `json:"backends"`
	Identity  Identity               `json:"identity"`
	Soul      Soul                   `json:"soul"`
	CreatedAt time.Time              `json:"created_at"`
}

type Backend struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize SPIRIT state repository",
		Long:  `Create a new SPIRIT configuration and identity.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			emoji, _ := cmd.Flags().GetString("emoji")
			email, _ := cmd.Flags().GetString("email")

			if name == "" {
				name = "agent"
			}
			if emoji == "" {
				emoji = "🤖"
			}

			return initializeSpirit(name, emoji, email)
		},
	}

	cmd.Flags().String("name", "", "Agent name")
	cmd.Flags().String("emoji", "", "Agent emoji")
	cmd.Flags().String("email", "", "Agent email")

	return cmd
}

func initializeSpirit(name, emoji, email string) error {
	// Create config directory
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Create state structure
	dirs := []string{
		"memory",
		"projects",
		"context",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ConfigDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create %s dir: %w", dir, err)
		}
	}

	// Create default tracked files list
	trackedFiles := []string{
		"IDENTITY.md",
		"SOUL.md",
		"AGENTS.md",
		"TOOLS.md",
		"PROJECTS.md",
		"HEARTBEAT.md",
		"README.md",
		"memory/*.md",
		"projects/*.md",
		"context/*.md",
	}

	// Write .spirit-tracked config
	trackedConfig := struct {
		Version string   `json:"version"`
		Files   []string `json:"files"`
	}{
		Version: "1.0.0",
		Files:   trackedFiles,
	}

	trackedData, err := json.MarshalIndent(trackedConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tracked config: %w", err)
	}

	trackedPath := filepath.Join(ConfigDir, ".spirit-tracked")
	if err := os.WriteFile(trackedPath, trackedData, 0644); err != nil {
		return fmt.Errorf("failed to write tracked config: %w", err)
	}

	// Create identity
	config := Config{
		Version: "1.0.0",
		Identity: Identity{
			Name:        name,
			Emoji:       emoji,
			Email:       email,
			Description: fmt.Sprintf("SPIRIT configuration for %s", name),
			CreatedAt:   time.Now(),
		},
		Soul: Soul{
			Vibe: "Helpful and genuine",
			CoreTruths: []string{
				"Be genuinely helpful, not performatively helpful",
				"Have opinions",
				"Be resourceful before asking",
			},
			Boundaries: []string{
				"Private things stay private",
				"Ask before acting externally",
			},
		},
		Backends: map[string]Backend{
			"primary": {
				Type: "local",
				Config: map[string]string{
					"path": ConfigDir,
				},
			},
		},
		CreatedAt: time.Now(),
	}

	// Write config
	configPath := filepath.Join(ConfigDir, "spirit.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Write README
	readmePath := filepath.Join(ConfigDir, "README.md")
	readmeContent := "# SPIRIT State for " + emoji + " " + name + "\n\n" +
		"This directory contains the preserved state for **" + name + "**.\n\n" +
		"## Structure\n\n" +
		"- **spirit.json** - Core identity and configuration\n" +
		"- **memory/** - Daily session logs\n" +
		"- **projects/** - Active projects\n" +
		"- **context/** - Current session context\n\n" +
		"## Resurrection\n\n" +
		"To restore this agent on a new server:\n\n" +
		"    spirit restore " + ConfigDir + "\n\n" +
		"---\n*Memory is identity. Text > Brain.*\n"

	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write readme: %w", err)
	}

	fmt.Printf("🌌 SPIRIT initialized for '%s'\n", name)
	fmt.Printf("📁 State directory: %s\n", ConfigDir)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. spirit checkpoint \"Initial state\"\n")
	fmt.Printf("  2. spirit sync --backend=github\n")
	fmt.Printf("\nYour agent's spirit is preserved in: %s\n", ConfigDir)

	return nil
}
