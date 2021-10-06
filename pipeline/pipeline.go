package main

import (
	"os"
	"path/filepath"

	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

type Transform interface {
	transform(input string) string
}

type Transformer struct {
	transform func(input string) string
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

func (p *Pipeline) New(workerCount int) *Pipeline {
	return &Pipeline{
		workers:      workerCount,
		transformers: []*Transformer{},
	}
}

func (m *TextMetadata) Run() interface{} {
	text := m.content
	for _, t := range m.transformers {
		text = t.transform(text)
	}
	return text
}

func (p *Pipeline) RegisterTransformers(transformerList []*Transformer) {
	p.transformers = transformerList
}

func (p *Pipeline) generateTransformWorkers() (chan concurrently.WorkFunction, <-chan concurrently.OrderedOutput) {
	inputCh := make(chan concurrently.WorkFunction)
	output := concurrently.Process(inputCh, &concurrently.Options{PoolSize: 10, OutChannelBuffer: 10})
	return inputCh, output
}

func (p *Pipeline) Execute(dir string) error {
	inputCh, output := p.generateTransformWorkers()

	err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			err := p.readAndSendFile(file, inputCh)
			if err != nil {
				return err
			}
		}
		return nil
	})

	for out := range output {
		go p.writeFile(out)
	}

	return err
}
