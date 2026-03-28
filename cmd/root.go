package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eitanoid/toml-resume/internal/resume"
	"github.com/spf13/cobra"
)

// flag values
var (
	outputFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "resume-gen [input file]",
	Short: "Generate a LaTeX resume from TOML/YAML",
	Args:  cobra.ExactArgs(1), // ensures exactly one input file
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		// resolve output file, infer input.tex if not specified
		finalOutput := outputFile
		if finalOutput == "" {
			ext := filepath.Ext(inputPath)
			finalOutput = strings.TrimSuffix(inputPath, ext) + ".tex"
		}

		resume := resume.NewDefaultResume()
		if err := resume.LoadDataFromFile(inputPath); err != nil {
			return err
		}
		resume.CreateLatexDoc()

		if finalOutput == "-" {
			// write to stdout
			fmt.Print(resume.String())
		} else {
			// write to file
			err := os.WriteFile(finalOutput, []byte(resume.String()), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("Successfully generated: %s\n", finalOutput)
		}

		return nil
	},
}

func init() {
	// Define the output flag. Default is empty string to trigger inference logic.
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output path (use '-' for stdout)")
}

// func main() {
// 	var inputPath string
// 	var outputPath string
//
// 	if len(os.Args) > 1 {
// 		inputPath = os.Args[1]
// 	}
//
// 	if inputPath == "" {
// 		fmt.Fprint(os.Stderr, ErrNoInput)
// 		os.Exit(1)
// 	}
//
// 	outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".tex"
//
// 	resume := NewDefaultResume()
// 	if err := resume.loadDataFromFile(inputPath); err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to load data from specified file: %v", err)
// 	}
// 	if err := resume.ValidateConfig(); err != nil {
// 		fmt.Fprintf(os.Stderr, "invalid configuration: %v", err)
// 	}
//
// 	resume.CreateLatexDoc()
//
// 	if outputPath != "" {
// 		f, err := os.Create(outputPath)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "failed to create output file \"%s\": %v", outputPath, err)
// 			os.Exit(1)
// 		}
// 		defer f.Close()
// 		if _, err := f.WriteString(resume.String()); err != nil {
// 			fmt.Fprintf(os.Stderr, "failed to write to output file \"%s\": %v", outputPath, err)
// 			os.Exit(1)
// 		}
// 	} else {
// 		fmt.Println(resume.String())
// 	}
//
// }
