package processors

import (
	"context"
	"os"
	"strings"
	"sync"

	astjson "GoOperatorAST/ast_json"
	"github.com/google/go-github/github"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

var logger *zap.Logger
var OUTPUT_DIR = os.Getenv("OUTPUT_DIR")
var fileProcessor *ants.PoolWithFunc

// SetupProcessing sets up the output directory and logger for processing.
// @param outputDir string
// @param logger *zap.Logger
func SetupProcessing(outputDir string, logger *zap.Logger) {
	OUTPUT_DIR = outputDir
	logger = logger
}

type FileInfo struct {
	Content  *string
	Path     *string
	FileName *string
}

type FileProcessor struct {
	FileInfo *FileInfo
	Wg       *sync.WaitGroup
	Logger   *zap.Logger
	Err      error
}

// CreateProcessFilePool creates a pool of goroutines to process files.
// It takes the pool size as an argument and returns a pointer to the ants.PoolWithFunc and an error (if any).
func CreateProcessFilePool(poolsize int) (*ants.PoolWithFunc, error) {
	p, err := ants.NewPoolWithFunc(poolsize, func(i interface{}) {
		// Type assert the input as FileProcessor
		pf := i.(FileProcessor)
		// Defer the Done() method of the WaitGroup to mark the completion of the goroutine
		defer pf.Wg.Done()
		// Call the parseFile method of the FileProcessor and assign the error to the Err field
		pf.Err = pf.parseFile()
		// Return to exit the goroutine
		return
	})
	if err != nil {
		return nil, err
	}
	// Assign the created pool to the global variable fileProcessor
	fileProcessor = p
	// Return the created pool and nil error
	return p, nil
}

// parseFile parses the file and converts it to JSON format.
func (pf *FileProcessor) parseFile() error {
	// Specify options for converting the file to JSON.
	options := astjson.Options{
		WithImports:    true,
		WithComments:   true,
		WithPositions:  true,
		WithReferences: true,
	}

	// Set the indentation string.
	indentStr := strings.Repeat(" ", 2)

	// Convert the file to JSON with the specified options and indentation.
	err := astjson.SourceToJSONWithContent(pf.FileInfo.Content, *pf.FileInfo.Path, *pf.FileInfo.FileName, indentStr, options)
	if err != nil {
		pf.Logger.Error("unable to convert file to json", zap.Error(err))
		return err
	}

	return nil
}

// fetchAndParseFile fetches and parses a Go file from a GitHub repository
// @param ctx context.Context
// @param client *github.Client
// @param owner string
// @param repo string
// @param path string
// @param tag string
func fetchAndParseFile(ctx context.Context, client *github.Client, owner, repo, path string, tag string) {
	// Fetch Golang code from GitHub repository
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		logger.Error("Error fetching repository content", zap.Error(err))
		return
	}

	// Get the content of the file
	content, err := fileContent.GetContent()
	if err != nil {
		logger.Error("Error getting file content", zap.Error(err))
		return
	}

	// Generate the file name
	fileName := strings.ReplaceAll(path, "/", "_") + ".json"

	// Generate the fully qualified file name
	fqfn := OUTPUT_DIR + "/" + tag + "/" + fileName

	// Create a FileInfo struct
	fileInfo := FileInfo{
		Content:  &content,
		Path:     &path,
		FileName: &fqfn,
	}

	var wg sync.WaitGroup

	// Create a FileProcessor struct
	process := FileProcessor{
		FileInfo: &fileInfo,
		Wg:       &wg,
		Logger:   logger,
		Err:      nil,
	}

	// Invoke the file processor in a goroutine
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		err := fileProcessor.Invoke(process)
		if err != nil {
			logger.Error("Unable to invoke file processor", zap.Error(err))
			wg.Done()
			return
		}
		wg.Wait()
		logger.Debug("Finished processing file:", zap.String("path", *process.FileInfo.Path), zap.String("filename", *process.FileInfo.FileName))
	}(process.Wg)
}

// ProcessRepo fetches and parses Go files from a GitHub repository
func ProcessRepo(ctx context.Context, client *github.Client, owner, repo, dir string, tag string) {
	// Create a stack to store directory paths
	var stack []string
	stack = append(stack, dir)

	// Set options for fetching repository content
	opts := &github.RepositoryContentGetOptions{}
	if tag != "" {
		opts.Ref = tag
	}

	// Iterate through the stack until it is empty
	for len(stack) > 0 {
		// Pop an item from the stack
		currentDir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Fetch contents of the current directory
		_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, currentDir, opts)
		if err != nil {
			logger.Error("Error fetching repository content", zap.Error(err))
			return
		}

		// Iterate through the directory content
		for _, content := range directoryContent {
			if *content.Type == "file" && strings.HasSuffix(*content.Path, ".go") {
				// Fetch and parse the Go file
				fetchAndParseFile(ctx, client, owner, repo, *content.Path, tag)
			} else if *content.Type == "dir" {
				// Push new directory into the stack
				stack = append(stack, *content.Path)
			}
		}
	}
}
