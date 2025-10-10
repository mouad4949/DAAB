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
	services       []string
	config         *Config
	ConfigMicro    *ConfigMicro
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

/******************************************************/
/************EntryPoint*******************************/
/****************************************************/

func (i *Initializer) Run() error {
	// Initialize config with defaults
	i.config = NewConfig()
	i.ConfigMicro = NewConfigMicro()

	// Step 1: Ask interactive questions (or use defaults)
	if err := i.gatherUserInput(); err != nil {
		return err
	}

	//Step 2:gather information based on teh project type
	if i.config.ProjectType == "monolith" {
		i.detectProjectMonolith()
	} else {
		i.DetectProjectMicroservice()

	}
	// Step 4: Validate configuration
	if err := i.validateConfig(); err != nil {
		return err
	}

	// Step 5: Save configuration
	if i.config.ProjectType == "microservice" {
		if err := i.saveConfigMicroservice(); err != nil {
			return err
		}
	} else {
		if err := i.saveConfig(); err != nil {
			return err
		}
	}

	return nil
}

/******************************************************/
/************AbstractUserInput************************/
/****************************************************/

func (i *Initializer) gatherUserInput() error {
	if i.nonInteractive {
		fmt.Println("⚙️  Using default configuration (non-interactive mode)...")
		i.setDefaults()
		return nil
	}

	fmt.Println("📝 Please answer a few questions about your project:")
	fmt.Println()

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

	if i.config.ProjectType == "microservice" {
		i.ConfigMicro.ProjectType = projectType
		// Project name microservice
		projectName, err := promptString("you microservice Project name", i.getDefaultProjectName())
		if err != nil {
			return err
		}
		i.ConfigMicro.ProjectName = projectName

		// Cloud provider
		cloudProvider, err := promptSelect(
			"Cloud provider",
			[]string{"aws", "gcp", "azure"},
			"aws",
		)
		if err != nil {
			return err
		}
		i.ConfigMicro.CloudProvider = cloudProvider
	} else {
		// Project name monolith
		projectName, err := promptString("your monolith Project name", i.getDefaultProjectName())
		if err != nil {
			return err
		}
		i.config.ProjectName = projectName
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
	}

	fmt.Println()
	return nil
}

/******************************************************/
/************detectProjectMonolith********************/
/****************************************************/

func (i *Initializer) detectProjectMonolith() error {
	fmt.Println("🔍 Detecting project type...")
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
	result, err := i.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}

	i.config.Language = result.Language
	i.config.Framework = result.Framework
	i.config.DetectedFiles = result.DetectedFiles
	port, err := promptInt("Application port", i.getDefaultPort())
	if err != nil {
		return err
	}
	i.config.Port = port
	fmt.Printf("   Language: %s\n", result.Language)
	if result.Framework != "" {
		fmt.Printf("   Framework: %s\n", result.Framework)
	}
	fmt.Println("project path:", i.projectPath)

	return nil
}

/******************************************************/
/************DetectProjectMicroservice****************/
/****************************************************/

func (i *Initializer) DetectSubfolders(root string) ([]string, error) {

	var subfolders []string

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subfolders = append(subfolders, filepath.Join(root, entry.Name()))
		}
	}

	return subfolders, nil
}

func (i *Initializer) DetectProjectMicroservice() error {
	fmt.Println("🔍 Detecting Microservices technologies stack...")

	folders, err := i.DetectSubfolders(i.projectPath)

	if err != nil {
		return fmt.Errorf("Error in detecting subfolders: ", err)
	}

	fmt.Printf("📁 Found %d subfolders:\n", len(folders))
	for _, folder := range folders {
		fmt.Println(" -", folder)
	}

	for _, folder := range folders {
		reserve := i.projectPath
		i.projectPath = folder
		i.detector.projectPath = folder

		result, err := i.detector.Detect()
		if err != nil {
			fmt.Printf("❌ Detection failed for %s: %v\n", folder, err)
			continue
		}

		fmt.Printf("✅ %s detected as %s (%s)\n", folder, result.Language, result.Framework)

		i.config.Language = result.Language
		i.config.Framework = result.Framework
		i.config.DetectedFiles = result.DetectedFiles

		port, err := promptInt("Application port", i.getDefaultPort())
		if err != nil {
			return err
		}
		i.config.Port = port
		i.services = append(i.services, folder)
		if err := i.saveConfig(); err != nil {
			fmt.Println("error in creating daab.yaml in this project:", folder)
			return err
		}
		i.projectPath = reserve
	}

	i.ConfigMicro.DetectedMicroservices = i.services
	fmt.Println()

	return nil
}

/******************************************************/
/************SetDefaults*******************************/
/****************************************************/

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

/******************************************************/
/************VALIDATE*******************************/
/****************************************************/

func (i *Initializer) validateConfig() error {
	if i.config.ProjectType == "monolith" {
		if i.config.ProjectName == "" {
			return fmt.Errorf("project name cannot be empty")
		}
		if i.config.Language == "" {
			return fmt.Errorf("language detection failed")
		}
		if i.config.Port <= 0 || i.config.Port > 65535 {
			return fmt.Errorf("invalid port number: %d", i.config.Port)
		}
	} else {
		if i.services == nil {
			return fmt.Errorf("no microservices detected inside of the folder")
		}
	}
	return nil
}

/******************************************************/
/************SAVE*******************************/
/****************************************************/

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

func (i *Initializer) saveConfigMicroservice() error {
	// Create .init directory
	initDir := filepath.Join(i.projectPath, ".init")
	if err := os.MkdirAll(initDir, 0755); err != nil {
		return fmt.Errorf("failed to create .init directory: %w", err)
	}

	// Generate YAML content
	data, err := yaml.Marshal(i.ConfigMicro)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configPath := filepath.Join(initDir, "daab.root.yaml")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
