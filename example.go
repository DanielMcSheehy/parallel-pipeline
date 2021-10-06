package main

import (
	"strings"

	"github.com/DanielMcSheehy/parallel-pipeline/pipeline"

	"github.com/spf13/cobra"
)

// example text transformations
func RemoveAllSmileyFaces() *pipeline.Transformer {
	return &pipeline.Transformer{
		Transform: func(input string) string {
			return strings.ReplaceAll(input, "😀", "")
		},
	}
}

func ReplaceSadWithHappy() *pipeline.Transformer {
	return &pipeline.Transformer{
		Transform: func(input string) string {
			return strings.ReplaceAll(input, "sad", "happy")
		},
	}
}

var workerCount = 3

func main() {
	var cmdTransform = &cobra.Command{
		Use:   "transform [directory to read] [output directory]",
		Short: "transform all files in a directory",
		Long:  `transform all files in a directory.`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			mainPipeline := pipeline.New(workerCount)
			mainPipeline.RegisterTransformers(
				RemoveAllSmileyFaces(),
				ReplaceSadWithHappy(),
			)
			mainPipeline.Execute(args[0], args[1])
		},
	}

	cmdTransform.Flags().IntVarP(&workerCount, "workers", "w", 3, "number of concurrent workers")

	var rootCmd = &cobra.Command{Use: "data-pipeline"}
	rootCmd.AddCommand(cmdTransform)
	rootCmd.Execute()
}
