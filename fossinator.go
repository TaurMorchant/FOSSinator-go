package main

import (
	"fmt"
	"fossinator/config"
	"fossinator/fs"
	"fossinator/processor"
	"fossinator/validator"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	err := config.Load()
	if err != nil {
		fmt.Println("Cannot load config file.", err)
		os.Exit(1)
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: "fossinator"}

	var transformCmd = &cobra.Command{
		Use: "transform",
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			fmtFlag, _ := cmd.Flags().GetBool("fmt")
			tidyFlag, _ := cmd.Flags().GetBool("tidy")
			transform(dir, fmtFlag, tidyFlag)
		},
	}
	transformCmd.Flags().StringP("dir", "d", "", "Directory to process")
	transformCmd.Flags().Bool("fmt", false, "Run 'go fmt' step")
	transformCmd.Flags().Bool("tidy", false, "Run 'go mod tidy' step")

	var validateCmd = &cobra.Command{
		Use: "validate",
		Run: func(cmd *cobra.Command, args []string) {
			dir := getDir(cmd)
			validate(dir)
		},
	}
	validateCmd.Flags().StringP("dir", "d", "", "Directory to process")

	rootCmd.AddCommand(transformCmd, validateCmd)
	_ = rootCmd.Execute()
}

func transform(dir string, fmtFlag, tidyFlag bool) {
	if _, err := fs.FindGoModFile(dir); err != nil {
		fmt.Printf("Directory '%s' is not a go module, cannot continue", dir)
		os.Exit(1)
	}
	fmt.Println("Directory to process: ", dir)

	if err := processor.UpdateImports(dir); err != nil {
		fmt.Println("Error during update imports:", err)
	}

	if err := processor.UpdateGoMod(dir); err != nil {
		fmt.Println("Error during update go.mod:", err)
	}

	if err := processor.AddConfigLoaderConfiguration(dir); err != nil {
		fmt.Println("Error during AddConfigLoaderConfiguration:", err)
	}

	if fmtFlag {
		processor.RunGoCommand(dir, "fmt", "./...")
	}

	if tidyFlag {
		processor.RunGoCommand(dir, "mod", "tidy")
	}
}

func validate(dir string) {
	validationMessages := validator.Validate(dir)

	if len(validationMessages) > 0 {
		fmt.Println("Validation completed with errors:")
		for _, msg := range validationMessages {
			fmt.Println(msg)
		}
	} else {
		fmt.Println("No validation errors")
	}
}

func getDir(cmd *cobra.Command) string {
	dirFlag, err := cmd.Flags().GetString("dir")
	if err != nil || len(dirFlag) == 0 {
		fmt.Println("Directory not specified. Current directory will be used")
		return "."
	} else {
		return dirFlag
	}
}
