package initcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type InitFlags struct {
	NonInteractive bool
	ProjectPath    string
}

func NewInitCommand() *cobra.Command {
	flags := &InitFlags{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize DAAB configuration for your project",
		Long: `Detect your project's language/framework and generate a daab.yaml configuration file.
This command will:
  - Auto-detect your project type (Node.js, Go, Python, etc.)
  - Ask interactive questions about your deployment preferences
  - Create a .init/daab.yaml configuration file`,
		Example: `  daab init
  daab init --non-interactive
  daab init --project-path /path/to/project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(flags)
		},
	}

	cmd.Flags().StringVar(&flags.ProjectPath, "project-path", ".", "Path to the project directory")

	return cmd
}

func runInit(flags *InitFlags) error {
	fmt.Println("🚀 Initializing DAAB for your project...")
	fmt.Println()

	// Create the initializer
	initializer := NewInitializer(flags.ProjectPath)

	// Run the initialization process
	if err := initializer.Run(); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ DAAB initialization complete!")
	fmt.Println("📝 Configuration saved to .init/daab.yaml")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review the generated .init/daab.yaml file")
	fmt.Println("  2. Run 'daab generate' to create deployment files")
	fmt.Println("  3. Run 'daab deploy' to deploy your application")

	return nil
}
