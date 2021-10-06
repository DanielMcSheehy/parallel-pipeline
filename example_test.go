package main_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	p "github.com/DanielMcSheehy/parallel-pipeline"
	"github.com/DanielMcSheehy/parallel-pipeline/pipeline"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "transform [directory to read] [output directory]",
	Short: "transform all files in a directory",
	Long:  `transform all files in a directory.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mainPipeline := pipeline.New(3)
		mainPipeline.RegisterTransformers(
			p.RemoveAllSmileyFaces(),
			p.ReplaceSadWithHappy(),
		)
		mainPipeline.Execute("example", ".")
	},
}

// Creates a example file, runs the example pipeline,
// and verifies the pipeline is correct
func TestPipeline(t *testing.T) {
	content := []byte("sad sad sad")
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()

	data, _ := os.ReadFile(dir + "/tmpfile")

	assert.Equal(t, strings.ContainsAny(string(data), "happy"), true)
}
