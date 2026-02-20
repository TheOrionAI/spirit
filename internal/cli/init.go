package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	Version   string             `json:"version"`
	Backends  map[string]Backend `json:"backends"`
	Identity  Identity           `json:"identity"`
	Soul      Soul               `json:"soul"`
	CreatedAt time.Time          `json:"created_at"`
}

type Backend struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type TrackedConfig struct {
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize SPIRIT state repository",
		Long: `Create a new SPIRIT configuration and identity.

Use --workspace DIR to enable "workspace mode": the .spirit-tracked 
config file is created in your workspace and symlinked to ~/.spirit/,
so you can easily edit which files to track.

Examples:
  spirit init --name="orion" --emoji="🌌"
  spirit init --workspace=/root/.openclaw/workspace
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			emoji, _ := cmd.Flags().GetString("emoji")
			email, _ := cmd.Flags().GetString("email")
			workspace, _ := cmd.Flags().GetString("workspace")

			if name == "" {
				name = "agent"
			}
			if emoji == "" {
				emoji = "🤖"
			}

			if workspace != "" {
				return initializeSpiritWorkspace(name, emoji, email, workspace)
			}
			return initializeSpirit(name, emoji, email)
		},
	}

	cmd.Flags().String("name", "", "Agent name")
	cmd.Flags().String("emoji", "", "Agent emoji")
	cmd.Flags().String("email", "", "Agent email")
	cmd.Flags().String("workspace", "", "Workspace directory (enables symlinked config mode)")

	return cmd
}

func initializeSpiritWorkspace(name, emoji, email, workspaceDir string) error {
	// Validate workspace path
	absWorkspace, err := filepath.Abs(workspaceDir)
	if err != nil {
		return fmt.Errorf("invalid workspace path: %w", err)
	}
	workspaceDir = absWorkspace

	// Create SPIRIT config directory
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Create subdirectories in workspace
	dirs := []string{"memory", "projects", "context"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(workspaceDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create %s dir: %w", dir, err)
		}
	}

	// Default tracked files (OpenClaw-friendly)
	trackedConfig := TrackedConfig{
		Version: "1.0.0",
		Files: []string{
			"IDENTITY.md",
			"SOUL.md",
			"AGENTS.md",
			"USER.md",
			"TOOLS.md",
			"PROJECTS.md",
			"HEARTBEAT.md",
			"README.md",
			"memory/*.md",
			"projects/*.md",
			"context/*.md",
		},
	}

	// Write .spirit-tracked in workspace
	trackedData, _ := json.MarshalIndent(trackedConfig, "", "  ")
	workspaceTrackedPath := filepath.Join(workspaceDir, ".spirit-tracked")
	if err := os.WriteFile(workspaceTrackedPath, trackedData, 0644); err != nil {
		return fmt.Errorf("failed to write tracked config: %w", err)
	}

	// Create symlink: ~/.spirit/.spirit-tracked -> workspace/.spirit-tracked
	spiritTrackedPath := filepath.Join(ConfigDir, ".spirit-tracked")
	os.Remove(spiritTrackedPath) // Remove if exists
	if err := os.Symlink(workspaceTrackedPath, spiritTrackedPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	// Create identity files in workspace
	identityContent := fmt.Sprintf("# %s %s\n\nName: %s\nEmoji: %s\n", emoji, name, name, emoji)
	os.WriteFile(filepath.Join(workspaceDir, "IDENTITY.md"), []byte(identityContent), 0644)

	soulContent := "# SOUL\n\nTODO: Define personality, behavior, boundaries\n"
	os.WriteFile(filepath.Join(workspaceDir, "SOUL.md"), []byte(soulContent), 0644)

	// Write spirit.json with workspace reference
	config := Config{
		Version: "1.1.0",
		Identity: Identity{Name: name, Emoji: emoji, Email: email, CreatedAt: time.Now()},
		Backends: map[string]Backend{
			"workspace": {Type: "workspace", Config: map[string]string{"path": workspaceDir}},
		},
		CreatedAt: time.Now(),
	}
	configData, _ := json.MarshalIndent(config, "", "  ")
	configPath := filepath.Join(ConfigDir, "spirit.json")
	os.WriteFile(configPath, configData, 0600)

	// Write README
	readmeContent := fmt.Sprintf(`# SPIRIT State for %s %s

## Structure
- **.spirit-tracked**: Edit this to control what syncs (symlinked from ~/.spirit/)
- **IDENTITY.md, SOUL.md**: Agent identity
- **memory/**, **projects/**, **context/**: State directories

## Sync
To sync with workspace source:
  SPIRIT_SOURCE_DIR=%s spirit sync
`, emoji, name, workspaceDir)
	os.WriteFile(filepath.Join(workspaceDir, "README.md"), []byte(readmeContent), 0644)

	fmt.Printf("🌌 SPIRIT initialized in workspace mode\n")
	fmt.Printf("📁 Workspace: %s\n", workspaceDir)
	fmt.Printf("🔗 Config symlink: ~/.spirit/.spirit-tracked -> %s\n", workspaceTrackedPath)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Edit %s/.spirit-tracked to configure files to sync\n", workspaceDir)
	fmt.Printf("2. Set SPIRIT_SOURCE_DIR=%s then run 'spirit sync'\n", workspaceDir)
	return nil
}

func initializeSpirit(name, emoji, email string) error {
	// Standard init (create in ~/.spirit/)
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	dirs := []string{"memory", "projects", "context"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ConfigDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create %s dir: %w", dir, err)
		}
	}

	trackedConfig := TrackedConfig{
		Version: "1.0.0",
		Files: []string{
			"IDENTITY.md", "SOUL.md", "AGENTS.md", "TOOLS.md",
			"memory/*.md", "projects/*.md", "context/*.md",
		},
	}
	trackedData, _ := json.MarshalIndent(trackedConfig, "", "  ")
	os.WriteFile(filepath.Join(ConfigDir, ".spirit-tracked"), trackedData, 0644)

	config := Config{
		Version: "1.0.0",
		Identity: Identity{Name: name, Emoji: emoji, Email: email, CreatedAt: time.Now()},
		CreatedAt: time.Now(),
	}
	configData, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(filepath.Join(ConfigDir, "spirit.json"), configData, 0600)

	readmeContent := fmt.Sprintf("# SPIRIT State for %s %s\n\nRun: spirit sync\n", emoji, name)
	os.WriteFile(filepath.Join(ConfigDir, "README.md"), []byte(readmeContent), 0644)

	fmt.Printf("🌌 SPIRIT initialized for '%s'\n", name)
	fmt.Printf("📁 State directory: %s\n", ConfigDir)
	return nil
}

func formatBulletList(items []string) string {
	var result strings.Builder
	for _, item := range items {
		result.WriteString(fmt.Sprintf("- %s\n", item))
	}
	return result.String()
}
