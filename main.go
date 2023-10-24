package main

import (
	"GoOperatorAST/fileinfo"
	"context"
	"flag"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"go/parser"
	"go/token"
	"golang.org/x/oauth2"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
var GITHUB_OWNER = os.Getenv("GITHUB_OWNER")
var OUTPUT_DIR = os.Getenv("OUTPUT_DIR")
var DEFAULT_STACK_SIZE = 1000000000
var REPO_OLD_TAG = ""
var REPO_NEW_TAG = ""
var logger *zap.Logger

func init() {
	// Initialize Zap logger
	l, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to initialize Zap logger")
	}
	logger = l

	defer logger.Sync() // Flushes buffer, if any

	// Generate random secret key
	if err != nil {
		logger.Fatal("Failed to generate random bytes for secret key", zap.Error(err))
	}
}

func fetchAndParseFile(ctx context.Context, client *github.Client, owner, repo, path string, tag string) {
	// Fetch Golang code from GitHub repository
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		logger.Error("Error fetching repository content", zap.Error(err))
		return
	}
	content, err := fileContent.GetContent()
	if err != nil {
		logger.Error("Error getting file content", zap.Error(err))
		return
	}

	// Parse the fetched Golang code into an AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, strings.NewReader(content), parser.ParseComments)
	if err != nil {
		logger.Error("Error parsing code", zap.Error(err))
		return
	}

	// Create a FlatBuffer builder
	builder := flatbuffers.NewBuilder(0)

	// Create a string and get its offset
	filename := builder.CreateString(path)

	// Calculate the number of imports in the parsed Go file
	numImports := len(node.Imports)

	// Use the generated `GoFileInfoStart` method to start building `GoFileInfo`
	fileinfo.GoFileInfoStart(builder)
	fileinfo.GoFileInfoAddFilename(builder, filename)
	fileinfo.GoFileInfoAddNumImports(builder, int32(numImports))
	fileInfo := fileinfo.GoFileInfoEnd(builder)

	// Finish serializing by specifying the root object and completing the FlatBuffer
	builder.Finish(fileInfo)

	// Get the serialized byte slice
	buf := builder.FinishedBytes()

	if tag == "" {
		tag = "main"
	}

	// Serialize and save the FlatBuffer into a file
	fileName := strings.ReplaceAll(path, "/", "_") + ".fb"
	file, err := os.Create(OUTPUT_DIR + "/" + tag + "/" + fileName)
	if err != nil {
		logger.Error("Error creating file", zap.Error(err))
		return
	}
	defer file.Close()

	_, err = file.Write(buf)
	if err != nil {
		logger.Error("Error writing FlatBuffer file", zap.Error(err))
		return
	}

	// Print for demonstration purposes
	logger.Debug(fmt.Sprintf("File info of %s has been serialized into %s\n", path, fileName))
}

func processRepo(ctx context.Context, client *github.Client, owner, repo, dir string, tag string) {
	var stack []string
	stack = append(stack, dir)

	opts := &github.RepositoryContentGetOptions{}
	if tag != "" {
		opts.Ref = tag
	}

	for len(stack) > 0 {
		// Pop an item from the stack
		currentDir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, currentDir, opts)
		if err != nil {
			logger.Error("Error fetching repository content", zap.Error(err))
			return
		}

		for _, content := range directoryContent {
			if *content.Type == "file" && strings.HasSuffix(*content.Path, ".go") {
				fetchAndParseFile(ctx, client, owner, repo, *content.Path, tag)
			} else if *content.Type == "dir" {
				// Push new directory into the stack
				stack = append(stack, *content.Path)
			}
		}
	}
}

func main() {
	// Define flags for repository and token
	repo := flag.String("repo", "", "Required: Repository name")
	repoOldTag := flag.String("tagOld", "", "Required: Tag name for previous version")
	repoNewTag := flag.String("tagNew", "", "Option: Tag name for new version,empty new tag will use main")
	githubToken := flag.String("token", "", "Optional: Access token")
	githubOwner := flag.String("owner", "", "Optional: Repo owner name")

	// Parse the command-line arguments
	flag.Parse()

	// Check if the required repository argument is provided
	if *repo == "" {
		logger.Info("Please provide the repository argument.")
		return
	}

	if *repoOldTag == "" {
		logger.Info("Please provide the previous tag version argument.")
		return
	}

	if *githubToken != "" {
		GITHUB_TOKEN = *githubToken
	}

	if *githubOwner != "" {
		GITHUB_OWNER = *githubOwner
	}

	if *repoOldTag != "" {
		REPO_OLD_TAG = *repoOldTag
	}

	if *repoNewTag != "" {
		REPO_NEW_TAG = *repoNewTag
	} else {
		logger.Info("Using main as the new tag")
	}

	if OUTPUT_DIR == "" {
		OUTPUT_DIR = "./output_temp"
	}

	debug.SetMaxStack(DEFAULT_STACK_SIZE * 1)

	// Initialize GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	//Check and Create the output directory
	createDir(OUTPUT_DIR)
	createDir(OUTPUT_DIR + "/" + REPO_OLD_TAG)
	if REPO_NEW_TAG == "" {
		createDir(OUTPUT_DIR + "/" + "main")
	} else {
		createDir(OUTPUT_DIR + "/" + REPO_NEW_TAG)
	}

	var wg sync.WaitGroup
	// Process the repository
	wg.Add(1)
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		logger.Debug(fmt.Sprintf("Processing old repository %s with tag", *repo, REPO_OLD_TAG))
		processRepo(ctx, client, GITHUB_OWNER, *repo, "", REPO_OLD_TAG)
	}(&wg)

	wg.Add(1)
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		logger.Debug(fmt.Sprintf("Processing new repository %s with tag", *repo, REPO_NEW_TAG))
		processRepo(ctx, client, GITHUB_OWNER, *repo, "", REPO_NEW_TAG)
	}(&wg)
	fmt.Println("Waiting for all goroutines to finish...")
	wg.Wait()
	fmt.Println("All goroutines finished.")
}

func createDir(directory string) {

	// Check if the directory already exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Create the directory recursively
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			logger.Error("Error creating directory", zap.Error(err))
			return
		}
		logger.Debug("Directory created successfully.")
	} else {
		logger.Debug("Directory already exists.")
	}
}
