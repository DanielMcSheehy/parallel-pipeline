package main

import (
	"fmt"
	"io"
	"log"
	"os"

	orderedconcurrently "github.com/tejzpr/ordered-concurrently"
	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

func (p *Pipeline) readAndSendFile(file os.FileInfo, inputCh chan orderedconcurrently.WorkFunction) error {
	size := file.Size()

	f, err := os.Open(file.Name())
	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return err
	}

	defer f.Close() // Do not forget to close the file

	// create buffer
	// max chunk size is number of works
	index := 0
	b := make([]byte, size/int64(p.workers))

	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		inputCh <- TextMetadata{
			fileName:     file.Name(),
			content:      string(b[:readTotal]),
			transformers: p.transformers,
		}
		index++
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) writeFile(output concurrently.OrderedOutput) error {
	metadata := output.Value.(TextMetadata)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(metadata.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte(metadata.content)); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
