package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/eitanoid/toml-resume/internal/loader"
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
	Args:  cobra.MaximumNArgs(1), // ensures exactly one input file
	RunE: func(cmd *cobra.Command, args []string) error {
		finalOutput := outputFile
		var r io.Reader

		// read from file if input was provided
		if len(args) > 0 && args[0] != "-" {
			inputPath := args[0]

			// resolve output file, infer input.tex if not specified
			if finalOutput == "" {
				ext := filepath.Ext(inputPath)
				finalOutput = strings.TrimSuffix(inputPath, ext) + ".tex"
			}

			// open file as a reader
			var err error
			if r, err = os.Open(inputPath); err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}

		} else {
			// determine if input is provided from pipe or redirection
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				r = os.Stdin
			} else {
				// no pipe and no file argument
				return fmt.Errorf("no input provided via pipe or file argument")
			}
		}

		resume := resume.NewDefaultResume()

		if err := loader.LoadFromReader(r, resume.Data); err != nil {
			return fmt.Errorf("failed to load config data: %w", err)
		}
		resume.CreateLatexDoc()

		if finalOutput == "-" || finalOutput == "" {
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
