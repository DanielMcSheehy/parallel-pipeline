package pipeline

import (
	"os"
	"path/filepath"

	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

type Transform interface {
	Transform(input string) string
}

type Transformer struct {
	Transform func(input string) string
}

type Pipeline struct {
	workers      int
	transformers []*Transformer
}

type TextMetadata struct {
	fileName string
	content  string
	// need transformers to run them in sequence
	transformers []*Transformer
}

func New(workerCount int) *Pipeline {
	return &Pipeline{
		workers:      workerCount,
		transformers: []*Transformer{},
	}
}

func (m TextMetadata) Run() interface{} {
	text := m.content
	for _, t := range m.transformers {
		text = t.Transform(text)
	}
	return TextMetadata{
		fileName: m.fileName,
		content:  text,
	}
}

func (p *Pipeline) RegisterTransformers(transformerList ...*Transformer) {
	p.transformers = transformerList
}

func (p *Pipeline) generateTransformWorkers() (chan concurrently.WorkFunction, <-chan concurrently.OrderedOutput) {
	inputCh := make(chan concurrently.WorkFunction)
	output := concurrently.Process(inputCh, &concurrently.Options{PoolSize: 10, OutChannelBuffer: 10})
	return inputCh, output
}

func (p *Pipeline) Execute(dir, outputDir string) error {
	inputCh, output := p.generateTransformWorkers()
	quit := make(chan bool)

	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if file != nil && !file.IsDir() {
			go p.readAndSendFile(dir, file, inputCh, quit)
		}
		return nil
	})

	for out := range output {
		p.writeFile(outputDir, out, quit)
	}

	return err
}
