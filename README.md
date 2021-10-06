# Parallel Pipline
A blazing fast library that allows data pipelines to work in parallel. This can traverse and transform extremely large text files (100GB or more) in seconds. 
## Usage

```go
import "github.com/DanielMcSheehy/parallel-pipeline/pipeline"
```
Add some text transformations
```go
// example text transformation
func RemoveAllSmileyFaces() *pipeline.Transformer {
	return &pipeline.Transformer{
		Transform: func(input string) string {
			return strings.ReplaceAll(input, "😀", "")
		},
	}
}
```

```go
func main() {
			mainPipeline := pipeline.New(workerCount)
			mainPipeline.RegisterTransformers(
				RemoveAllSmileyFaces(),
				ReplaceSadWithHappy(),
			)
			mainPipeline.Execute(directory, ouputDirectory)
}
```