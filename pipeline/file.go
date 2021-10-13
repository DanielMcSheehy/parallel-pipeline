package pipeline

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

func (p *Pipeline) readAndSendFile(dir string, file os.FileInfo, inputCh chan concurrently.WorkFunction, quit chan bool) error {
	size := file.Size()

	f, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		fmt.Println("not able to read the file", err)
		return err
	}

	// create buffer
	// max chunk size is number of workers
	b := make([]byte, size/int64(p.workers)+1)

	defer f.Close() // Do not forget to close the file

	var wg sync.WaitGroup
	wg.Add(int(size / int64(p.workers)))

	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		fmt.Println("file: ", string(b[:readTotal]))

		inputCh <- TextMetadata{
			fileName:     file.Name(),
			content:      string(b[:readTotal]),
			transformers: p.transformers,
		}
		defer wg.Done()
	}

	close(inputCh)
	wg.Wait()

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
