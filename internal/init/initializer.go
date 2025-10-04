package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Initializer struct {
	projectPath    string
	nonInteractive bool
	config         *Config
	detector       *Detector
}

func NewInitializer(projectPath string, nonInteractive bool) *Initializer {
	return &Initializer{
		projectPath:    projectPath,
		nonInteractive: nonInteractive,
		config:         &Config{},
		detector:       NewDetector(projectPath),
	}
}

func (i *Initializer) Run() error {
	// Initialize config with defaults
	i.config = NewConfig()

	// Step 1: Detect project type
	if err := i.detectProject(); err != nil {
		return err
	}

	// Step 2: Ask interactive questions (or use defaults)
	if err := i.gatherUserInput(); err != nil {
		return err
	}

	// Step 3: Validate configuration
	if err := i.validateConfig(); err != nil {
		return err
	}

	// Step 4: Save configuration
	if err := i.saveConfig(); err != nil {
		return err
	}

	return nil
}

func (i *Initializer) detectProject() error {
	fmt.Println("🔍 Detecting project type...")

	result, err := i.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}

	i.config.Language = result.Language
	i.config.Framework = result.Framework
	i.config.DetectedFiles = result.DetectedFiles

	fmt.Printf("   Language: %s\n", result.Language)
	if result.Framework != "" {
		fmt.Printf("   Framework: %s\n", result.Framework)
	}
	fmt.Println()

	return nil
}

func (i *Initializer) gatherUserInput() error {
	if i.nonInteractive {
		fmt.Println("⚙️  Using default configuration (non-interactive mode)...")
		i.setDefaults()
		return nil
	}

	fmt.Println("📝 Please answer a few questions about your project:")
	fmt.Println()

	// Project name
	projectName, err := promptString("Project name", i.getDefaultProjectName())
	if err != nil {
		return err
	}
	i.config.ProjectName = projectName

	// Project type
	projectType, err := promptSelect(
		"Project type",
		[]string{"monolith", "microservice"},
		"monolith",
	)
	if err != nil {
		return err
	}
	i.config.ProjectType = projectType

	// Cloud provider
	cloudProvider, err := promptSelect(
		"Cloud provider",
		[]string{"aws", "gcp", "azure"},
		"aws",
	)
	if err != nil {
		return err
	}
	i.config.CloudProvider = cloudProvider

	// Port
	port, err := promptInt("Application port", i.getDefaultPort())
	if err != nil {
		return err
	}
	i.config.Port = port

	// Environment
	environment, err := promptString("Environment", "production")
	if err != nil {
		return err
	}
	i.config.Environment = environment

	// Container registry
	registry, err := promptString("Container registry (leave empty for default)", "")
	if err != nil {
		return err
	}
	i.config.ContainerRegistry = registry

	// Kubernetes namespace
	namespace, err := promptString("Kubernetes namespace", "default")
	if err != nil {
		return err
	}
	i.config.Namespace = namespace

	fmt.Println()
	return nil
}

func (i *Initializer) setDefaults() {
	i.config.ProjectName = i.getDefaultProjectName()
	i.config.ProjectType = "monolith"
	i.config.CloudProvider = "aws"
	i.config.Port = i.getDefaultPort()
	i.config.Environment = "production"
	i.config.Namespace = "default"
	i.config.ContainerRegistry = ""
}

func (i *Initializer) getDefaultProjectName() string {
	absPath, err := filepath.Abs(i.projectPath)
	if err != nil {
		return "my-app"
	}
	return filepath.Base(absPath)
}

func (i *Initializer) getDefaultPort() int {
	// Default ports based on language/framework
	portMap := map[string]int{
		"nodejs": 3000,
		"go":     8080,
		"python": 8000,
		"java":   8080,
		"ruby":   3000,
		"php":    8080,
		"dotnet": 5000,
		"rust":   8080,
	}

	if port, ok := portMap[i.config.Language]; ok {
		return port
	}
	return 8080
}

func (i *Initializer) validateConfig() error {
	if i.config.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if i.config.Language == "" {
		return fmt.Errorf("language detection failed")
	}
	if i.config.Port <= 0 || i.config.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", i.config.Port)
	}
	return nil
}

func (i *Initializer) saveConfig() error {
	// Create .init directory
	initDir := filepath.Join(i.projectPath, ".init")
	if err := os.MkdirAll(initDir, 0755); err != nil {
		return fmt.Errorf("failed to create .init directory: %w", err)
	}

	// Generate YAML content
	data, err := yaml.Marshal(i.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configPath := filepath.Join(initDir, "daab.yaml")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
