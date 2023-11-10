package main

import (
	"GoOperatorAST/processors"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"os"
	"sync"
	"time"
)

var (
	GITHUB_TOKEN  = os.Getenv("GITHUB_TOKEN")
	GITHUB_OWNER  = os.Getenv("GITHUB_OWNER")
	OUTPUT_DIR    = os.Getenv("OUTPUT_DIR")
	REPO_OLD_TAG  = ""
	REPO_NEW_TAG  = ""
	logger        *zap.Logger
	fileProcessor *ants.PoolWithFunc
)

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

	fp, err := processors.CreateProcessFilePool(200)
	if err != nil {
		logger.Fatal("Failed to create file processor pool", zap.Error(err))
	}
	fileProcessor = fp
}

// Main function
func main() {
	defer fileProcessor.Release()

	// Define flags for repository and token
	repo := flag.String("repo", "", "Required: Repository name")
	repoOldTag := flag.String("tagOld", "", "Required: Tag name for previous version")
	repoNewTag := flag.String("tagNew", "", "Option: Tag name for new version, empty new tag will use main")
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

	// Initialize GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	//client := github.NewClient(tc)
	_ = github.NewClient(tc)

	// Check and Create the output directory
	createDir(OUTPUT_DIR)
	createDir(OUTPUT_DIR + "/" + REPO_OLD_TAG)

	if REPO_NEW_TAG == "" {
		createDir(OUTPUT_DIR + "/" + "main")
	} else {
		createDir(OUTPUT_DIR + "/" + REPO_NEW_TAG)
	}

	processors.SetupProcessing(OUTPUT_DIR, logger)

	// Start the worker reporter
	worker_ch := make(chan bool)
	go func() {
		for {
			logger.Debug("file processing worker threads stats", zap.Int("Running", fileProcessor.Running()), zap.Int("Waiting", fileProcessor.Waiting()), zap.Int("Free", fileProcessor.Free()))

			// Wait for and print values received from channels
			select {
			case <-worker_ch:
				logger.Info("stopping worker reporter")
				return
			default:
			}
			time.Sleep(5 * time.Second)
		}
	}()

	var wg sync.WaitGroup

	// Process the repository
	wg.Add(1)
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		logger.Debug(fmt.Sprintf("Processing old repository %s with tag", *repo, REPO_OLD_TAG))
		//processors.ProcessRepo(ctx, client, GITHUB_OWNER, *repo, "", REPO_OLD_TAG)
		processors.ReadFiles(OUTPUT_DIR + "/" + REPO_OLD_TAG)
	}(&wg)

	wg.Add(1)
	go func(wg1 *sync.WaitGroup) {
		defer wg1.Done()
		logger.Debug(fmt.Sprintf("Processing new repository %s with tag", *repo, REPO_NEW_TAG))
		//processors.ProcessRepo(ctx, client, GITHUB_OWNER, *repo, "", REPO_NEW_TAG)
		processors.ReadFiles(OUTPUT_DIR + "/" + REPO_NEW_TAG)
	}(&wg)

	logger.Info("Waiting for all file processing to finish...")
	wg.Wait()

	worker_ch <- true
	logger.Info("All file processing has completed...")
}

// Create directory if it does not exist
// @param directory: directory path
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
