package initcmd

import (
	"fmt"
	config "github.com/mouad4949/DAAB/internal/init/config"
	configMicroservice "github.com/mouad4949/DAAB/internal/init/config/microservice"
	configMonolith "github.com/mouad4949/DAAB/internal/init/config/monolith"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Initializer struct {
	projectPath string

	services []string

	//Used just to detect port in getDefaultPort() by it's language
	baseconfigapp *config.BaseConfigApp

	//Used to create a daab.yaml file for monolith projects
	configmonolith *configMonolith.ConfigMonolith

	//Used to create a daab.root.yaml file for microservices projects in the root of the folder
	ConfigMicroRoot *configMicroservice.ConfigMicroRoot

	//Used to create a daab.yaml file for microservices projects on each microservice
	ConfigMicro *configMicroservice.ConfigMicroservice
	detector    *Detector
}

func NewInitializer(projectPath string) *Initializer {
	return &Initializer{
		projectPath: projectPath,
		detector:    NewDetector(projectPath),
	}
}

/******************************************************/
/************EntryPoint*******************************/
/****************************************************/

func (i *Initializer) Run() error {
	// Initialize config with defaults
	i.configmonolith = configMonolith.NewConfigMonolith()
	i.ConfigMicroRoot = configMicroservice.NewConfigMicroRoot()
	i.ConfigMicro = configMicroservice.NewConfigMicroservice()
	i.baseconfigapp = config.NewBaseConfigApp()

	// Step 1: Ask interactive questions (or use defaults)
	if err := i.gatherUserInput(); err != nil {
		return err
	}

	//Step 2:gather information based on the project type
	if i.configmonolith.ProjectType == "monolith" {
		i.detectProjectMonolith()
	} else {
		i.DetectProjectMicroservice()
	}
	// Step 4: Validate configuration
	if err := i.validateConfig(); err != nil {
		return err
	}

	// Step 5: Save configuration
	if i.ConfigMicroRoot.ProjectType == "microservice" {
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

	if projectType == "microservice" {
		//Project Type
		i.ConfigMicroRoot.ProjectType = projectType
		// Project name microservice
		projectName, err := promptString("your microservices Project name", i.getDefaultProjectName())
		if err != nil {
			return err
		}
		i.ConfigMicroRoot.ProjectName = projectName

		// Cloud provider
		cloudProvider, err := promptSelect(
			"Cloud provider",
			[]string{"aws", "gcp", "azure"},
			"aws",
		)
		if err != nil {
			return err
		}
		i.ConfigMicroRoot.CloudProvider = cloudProvider

		// Environment
		environment, err := promptString("Environment", "production")
		if err != nil {
			return err
		}
		i.ConfigMicroRoot.Environment = environment

		//region
		region, err := promptString("Region", "")
		if err != nil {
			return err
		}
		i.ConfigMicroRoot.Region = region

		//namespace
		namespace, err := promptString("Kubernetes namespace", "default")
		if err != nil {
			return err
		}
		i.ConfigMicroRoot.Namespace = namespace

	} else {
		//Project Type
		i.configmonolith.ProjectType = projectType
		// Project name monolith
		projectName, err := promptString("your monolith Project name", i.getDefaultProjectName())
		if err != nil {
			return err
		}
		i.configmonolith.ProjectName = projectName
		// Cloud provider
		cloudProvider, err := promptSelect(
			"Cloud provider",
			[]string{"aws", "gcp", "azure"},
			"aws",
		)
		if err != nil {
			return err
		}
		i.configmonolith.CloudProvider = cloudProvider

		// Environment
		environment, err := promptString("Environment", "production")
		if err != nil {
			return err
		}
		i.configmonolith.Environment = environment

		//region
		region, err := promptString("Region", "")
		if err != nil {
			return err
		}
		i.configmonolith.Region = region

		// Container registry
		registry, err := promptString("Container registry (leave empty for default)", "")
		if err != nil {
			return err
		}
		i.configmonolith.ContainerRegistry = registry

		// Kubernetes namespace
		namespace, err := promptString("Kubernetes namespace", "default")
		if err != nil {
			return err
		}
		i.configmonolith.Namespace = namespace

	}

	fmt.Println()
	return nil
}

/******************************************************/
/************detectProjectMonolith********************/
/****************************************************/

func (i *Initializer) detectProjectMonolith() error {
	fmt.Println("🔍 Detecting project type...")
	result, err := i.detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}

	i.configmonolith.Language = result.Language
	i.configmonolith.Framework = result.Framework
	i.configmonolith.DetectedFiles = result.DetectedFiles
	i.baseconfigapp.Language = i.configmonolith.Language
	port, err := promptInt("Application port", i.getDefaultPort())
	if err != nil {
		return err
	}
	i.configmonolith.Port = port
	fmt.Printf("   Language: %s\n", result.Language)
	if result.Framework != "" {
		fmt.Printf("   Framework: %s\n", result.Framework)
	}

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
		if strings.HasPrefix(folder, ".") {
			continue
		}
		i.services = append(i.services, folder)
		reserve := i.projectPath
		i.projectPath = folder
		i.detector.projectPath = folder

		result, err := i.detector.Detect()
		if err != nil {
			fmt.Printf("❌ Detection failed for %s: %v\n", folder, err)
			continue
		}

		fmt.Printf("✅ %s detected as %s (%s)\n", folder, result.Language, result.Framework)
		i.ConfigMicro.ProjectName = folder
		i.ConfigMicro.ProjectType = i.ConfigMicroRoot.ProjectType
		i.ConfigMicro.Language = result.Language
		i.ConfigMicro.Framework = result.Framework
		i.ConfigMicro.DetectedFiles = result.DetectedFiles
		i.baseconfigapp.Language = i.ConfigMicro.Language
		port, err := promptInt("Application port", i.getDefaultPort())
		if err != nil {
			return err
		}

		i.ConfigMicro.Port = port

		//cloud
		i.ConfigMicro.CloudProvider = i.ConfigMicroRoot.CloudProvider

		// Container registry
		registry, err := promptString("Container registry (leave empty for default)", "")
		if err != nil {
			return err
		}
		i.ConfigMicro.ContainerRegistry = registry

		if err := i.saveConfig(); err != nil {
			fmt.Println("error in creating daab.yaml in this project:", folder)
			return err
		}
		i.projectPath = reserve
	}

	i.ConfigMicroRoot.DetectedMicroservices = i.services

	fmt.Println()

	return nil
}

/******************************************************/
/************SetDefaults*******************************/
/****************************************************/

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

	if port, ok := portMap[i.baseconfigapp.Language]; ok {
		return port
	}
	return 8080
}

/******************************************************/
/************VALIDATE*******************************/
/****************************************************/

func (i *Initializer) validateConfig() error {
	if i.configmonolith.ProjectType == "monolith" {
		if i.configmonolith.ProjectName == "" {
			return fmt.Errorf("project name cannot be empty")
		}
		if i.configmonolith.Language == "" {
			return fmt.Errorf("language detection failed")
		}
		if i.configmonolith.Port <= 0 || i.configmonolith.Port > 65535 {
			return fmt.Errorf("invalid port number: %d", i.configmonolith.Port)
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
	data, err := yaml.Marshal(i)
	if i.ConfigMicro.ProjectType == "microservice" {
		data, err = yaml.Marshal(i.ConfigMicro)
	} else {
		data, err = yaml.Marshal(i.configmonolith)
	}

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
	data, err := yaml.Marshal(i.ConfigMicroRoot)
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
