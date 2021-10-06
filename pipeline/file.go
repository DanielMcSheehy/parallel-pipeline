package pipeline

import (
	"fmt"
	"io"
	"log"
	"os"

	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

func (p *Pipeline) readAndSendFile(dir string, file os.FileInfo, inputCh chan concurrently.WorkFunction) error {
	size := file.Size()

	f, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		fmt.Println("not able to read the file", err)
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

func (p *Pipeline) writeFile(dir string, output concurrently.OrderedOutput) error {
	metadata := output.Value.(TextMetadata)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(dir+"/"+metadata.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(metadata.content)); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
