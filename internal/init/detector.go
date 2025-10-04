package initcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DetectionResult struct {
	Language      string
	Framework     string
	DetectedFiles []string
}

type Detector struct {
	projectPath string
}

func NewDetector(projectPath string) *Detector {
	return &Detector{
		projectPath: projectPath,
	}
}

func (d *Detector) Detect() (*DetectionResult, error) {
	result := &DetectionResult{
		DetectedFiles: []string{},
	}

	// Detection order matters - more specific checks first
	detectors := []func(*DetectionResult) bool{
		d.detectGo,
		d.detectNodeJS,
		d.detectPython,
		d.detectJava,
		d.detectRuby,
		d.detectPHP,
		d.detectDotNet,
		d.detectRust,
	}

	for _, detector := range detectors {
		if detector(result) {
			return result, nil
		}
	}

	return nil, fmt.Errorf("could not detect project language. Please ensure you're in a valid project directory")
}

func (d *Detector) detectNodeJS(result *DetectionResult) bool {
	packageJSON := filepath.Join(d.projectPath, "package.json")
	if !d.fileExists(packageJSON) {
		return false
	}

	result.Language = "nodejs"
	result.DetectedFiles = append(result.DetectedFiles, "package.json")

	// Try to detect framework
	data, err := os.ReadFile(packageJSON)
	if err == nil {
		var pkg map[string]interface{}
		if json.Unmarshal(data, &pkg) == nil {
			if deps, ok := pkg["dependencies"].(map[string]interface{}); ok {
				switch {
				case d.hasDep(deps, "express"):
					result.Framework = "express"
				case d.hasDep(deps, "next"):
					result.Framework = "nextjs"
				case d.hasDep(deps, "react"):
					result.Framework = "react"
				case d.hasDep(deps, "vue"):
					result.Framework = "vue"
				case d.hasDep(deps, "nestjs"):
					result.Framework = "nestjs"
				}
			}
		}
	}

	return true
}

func (d *Detector) detectGo(result *DetectionResult) bool {
	goMod := filepath.Join(d.projectPath, "go.mod")
	if !d.fileExists(goMod) {
		return false
	}

	result.Language = "go"
	result.DetectedFiles = append(result.DetectedFiles, "go.mod")

	// Try to detect framework by reading go.mod
	data, err := os.ReadFile(goMod)
	if err == nil {
		content := string(data)
		switch {
		case strings.Contains(content, "github.com/gin-gonic/gin"):
			result.Framework = "gin"
		case strings.Contains(content, "github.com/gofiber/fiber"):
			result.Framework = "fiber"
		case strings.Contains(content, "github.com/labstack/echo"):
			result.Framework = "echo"
		case strings.Contains(content, "github.com/gorilla/mux"):
			result.Framework = "gorilla"
		}
	}

	return true
}

func (d *Detector) detectPython(result *DetectionResult) bool {
	files := []string{"requirements.txt", "Pipfile", "pyproject.toml", "setup.py"}

	for _, file := range files {
		if d.fileExists(filepath.Join(d.projectPath, file)) {
			result.Language = "python"
			result.DetectedFiles = append(result.DetectedFiles, file)

			// Try to detect framework
			if file == "requirements.txt" {
				data, err := os.ReadFile(filepath.Join(d.projectPath, file))
				if err == nil {
					content := strings.ToLower(string(data))
					switch {
					case strings.Contains(content, "flask"):
						result.Framework = "flask"
					case strings.Contains(content, "django"):
						result.Framework = "django"
					case strings.Contains(content, "fastapi"):
						result.Framework = "fastapi"
					}
				}
			}
			return true
		}
	}
	return false
}

func (d *Detector) detectJava(result *DetectionResult) bool {
	if d.fileExists(filepath.Join(d.projectPath, "pom.xml")) {
		result.Language = "java"
		result.Framework = "maven"
		result.DetectedFiles = append(result.DetectedFiles, "pom.xml")
		return true
	}

	if d.fileExists(filepath.Join(d.projectPath, "build.gradle")) ||
		d.fileExists(filepath.Join(d.projectPath, "build.gradle.kts")) {
		result.Language = "java"
		result.Framework = "gradle"
		result.DetectedFiles = append(result.DetectedFiles, "build.gradle")
		return true
	}

	return false
}

func (d *Detector) detectRuby(result *DetectionResult) bool {
	if d.fileExists(filepath.Join(d.projectPath, "Gemfile")) {
		result.Language = "ruby"
		result.DetectedFiles = append(result.DetectedFiles, "Gemfile")

		// Check for Rails
		data, err := os.ReadFile(filepath.Join(d.projectPath, "Gemfile"))
		if err == nil && strings.Contains(string(data), "rails") {
			result.Framework = "rails"
		}
		return true
	}
	return false
}

func (d *Detector) detectPHP(result *DetectionResult) bool {
	if d.fileExists(filepath.Join(d.projectPath, "composer.json")) {
		result.Language = "php"
		result.DetectedFiles = append(result.DetectedFiles, "composer.json")

		// Try to detect framework
		data, err := os.ReadFile(filepath.Join(d.projectPath, "composer.json"))
		if err == nil {
			content := strings.ToLower(string(data))
			switch {
			case strings.Contains(content, "laravel"):
				result.Framework = "laravel"
			case strings.Contains(content, "symfony"):
				result.Framework = "symfony"
			}
		}
		return true
	}
	return false
}

func (d *Detector) detectDotNet(result *DetectionResult) bool {
	files := []string{"*.csproj", "*.fsproj", "*.vbproj"}

	for _, pattern := range files {
		matches, _ := filepath.Glob(filepath.Join(d.projectPath, pattern))
		if len(matches) > 0 {
			result.Language = "dotnet"
			result.DetectedFiles = append(result.DetectedFiles, filepath.Base(matches[0]))
			return true
		}
	}
	return false
}

func (d *Detector) detectRust(result *DetectionResult) bool {
	if d.fileExists(filepath.Join(d.projectPath, "Cargo.toml")) {
		result.Language = "rust"
		result.DetectedFiles = append(result.DetectedFiles, "Cargo.toml")
		return true
	}
	return false
}

// Helper functions
func (d *Detector) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (d *Detector) hasDep(deps map[string]interface{}, name string) bool {
	_, ok := deps[name]
	return ok
}
