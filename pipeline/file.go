package pipeline

import (
	"bufio"
	"fmt"
	"log"
	"os"

	concurrently "github.com/tejzpr/ordered-concurrently/v2"
)

func (p *Pipeline) readAndSendFile(dir string, file os.FileInfo, inputCh chan concurrently.WorkFunction, quit chan bool) error {
	f, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		fmt.Println("not able to read the file", err)
		return err
	}

	// Do not forget to close the file
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Not actually needed since it’s a default split function.
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println("file: ", scanner.Text())
		inputCh <- TextMetadata{
			fileName:     file.Name(),
			content:      scanner.Text(),
			transformers: p.transformers,
		}
	}

	close(inputCh)

	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) writeFile(dir string, output concurrently.OrderedOutput, quit chan bool) error {
	metadata := output.Value.(TextMetadata)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(dir+"/"+metadata.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(metadata.content + "\n")); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
