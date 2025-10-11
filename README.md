
# DAAB INIT - Project Initialization Module
This document provides a technical overview of the daab init command implementation.  
It explains the internal structure, key functions, and overall workflow so contributors can quickly understand and navigate the codebase.

## Overview

The daab init command initializes a project for DAAB deployment by:


- Detecting language/framework automatically



- Prompting the user for configuration details (interactive mode)



- Validating configuration data



- Saving .init/daab.yaml or .init/daab.root.yaml depending on project
  type

## Project Structure

cmd/  
└── daab/  
└── main.go

internal/  
└── init/  
├── command.go  
├── initializer.go  
├── detector.go  
├── prompts.go  
└── config/  
├── baseconfig.go  
├── baseappconfig.go  
├── monolith/  
│ └── configMonolith.go  
└── microservice/  
├── configMicro.go  
└── configMicroRoot.go


## Quick Reference

| File | Purpose |  
|--|--|
|**cmd/daab/main.go**  | CLI entry point; sets up Cobra and routes commands|
| **internal/init/command.go** | Registers the `init` command and handles its execution |
| **internal/init/initializer.go** | Core orchestration logic for project initialization|
| **internal/init/detector.go**| Detects language and framework automatically|
|**internal/init/prompts.go**|Handles user prompts (interactive mode)  |
| **internal/init/config/baseconfig.go**| Base configuration abstraction|
| **internal/init/config/baseappconfig.go**| Shared structure for monolith/micro configs|
| **internal/init/config/monolith/configMonolith.go**| Monolith-specific YAML structure |
| **internal/init/config/microservice/configMicro.go**|Microservice-specific YAML structure |
| **internal/init/config/microservice/configMicroRoot.go**| Root config for microservice projects |


## Key Functions by File


### cmd/daab/main.go
|Function  | Description |
|--|--|
| `main()` | CLI entry point; initializes Cobra and registers all commands |

### **internal/init/initializer.go**

|Function  | Description |
|--|--|
|`NewInitializer()`  |  Creates new initializer instance|
|`Run()`  |  Executes full init flow: detect → ask → validate → save|
| `detectProjectMonolith()` | Detects monolith project language and framework|
| `DetectProjectMicroservice()`| Detects each microservice language/framework|
|`DetectSubfolders()`|Detects number of microservices in root  |
|`gatherUserInput()`| Collects user input via prompt |
| `setDefaults()`| Sets default values for non-interactive mode |
|`getDefaultProjectName()` | Extracts project name from folder |
|`getDefaultPort()`| Suggests port based on language|
| `validateConfig()`| Validates config before saving|
|`saveConfig()`|Creates `.init/daab.yaml` for monoliths|
| `saveConfigMicroservice()`|Creates `.init/daab.root.yaml` for microservice roots|

### **internal/init/detector.go**

| Function | Description |
|--|--|
| `NewDetector()`| Creates new language detector instance|
|`Detect()`| Runs all detectors sequentially|
| `detectNodeJS()`| Detects Node.js (Express, Next, etc.)|
| `detectGo()`| Detects Go (Gin, Fiber, etc.)|
| `detectPython()`|Detects Python (Flask, FastAPI, etc.)|
| `detectJava()`|Detects Java (Maven/Gradle)|
| `detectRuby()`| Detects Ruby (Rails)|
| `detectPHP()`| Detects PHP (Laravel/Symfony)|
| `detectDotNet()` | Detects .NET projects|
| `detectRust()`| Detects Rust (Cargo.toml)|
| `fileExists()` |Checks if a file exists|
| `hasDep()`| Checks dependency existence in package.json|



### **internal/init/config/**

| File | Purpose |
|--|--|
| **baseconfig.go**| Base configuration data structure|
| **baseappconfig.go**| Shared abstract for app-level configs |
| **monolith/configMonolith.go**| Defines `daab.yaml` schema for monoliths|
| **microservice/configMicro.go**| Defines per-service YAML schema|
| **microservice/configMicroRoot.go**|  Defines root YAML schema (`daab.root.yaml`)|

### **internal/init/prompts.go**













|  Function|Description|
|--|--|
| `promptString()`| Prompts text input with default|
| `promptInt()`| Prompts numeric input with default|
| `promptSelect()`| Prompts selection from list|
|`promptConfirm()`| Prompts yes/no confirmation|


### Execution Flow Summary

1.  **CLI Entry** → User runs `daab init`

2.  **Command Setup** → Cobra triggers `runInit()`

3.  **Initializer** → Creates config handler via `NewInitializer()`

4.  **Detection Phase** → Automatically detects project type and language

5.  **User Input** → Interactive or non-interactive configuration

6.  **Validation** → Ensures config integrity

7.  **Saving** → Writes `.init/daab.yaml` (monolith) or `.init/daab.root.yaml` (microservices)

8.  **Completion** → CLI confirms initialization success

## 💡 Notes for Contributors

-   Keep detection logic modular (easy to add new languages).

-   Ensure YAML schema consistency between monolith and microservice configs.

-   Maintain prompt defaults aligned with detection results.



**Maintainer:** Mouad Rguibi  
**Language:** Go  
**Command:** `daab init`  
**Goal:** Zero-friction initialization and configuration generation for cloud deployment





