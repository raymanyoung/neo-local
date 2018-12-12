package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	verbose = false
)

func main() {
	travisCommitRange := os.Getenv("TRAVIS_COMMIT_RANGE")
	if travisCommitRange == "" {
		log.Fatalf("Expecting $TRAVIS_COMMIT_RANGE to not be empty")
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "-v" {
		verbose = true
	}

	filesChangedBytes, err := exec.Command(
		"git", "--no-pager", "diff", "--name-only", travisCommitRange,
	).Output()
	if err != nil {
		log.Printf("Files changed output: '%s'", string(filesChangedBytes))
		log.Fatal(err)
	}

	filesChanged := string(filesChangedBytes)
	filesChanged = strings.TrimSuffix(filesChanged, "\n")

	if filesChanged == "" {
		fmt.Println("none")
		os.Exit(0)
	}

	directories := make(map[string]int)

	logWithVerbose("Files changed:")
	for _, file := range strings.Split(filesChanged, "\n") {
		logfWithVerbose(" - %s", file)
		if _, ok := directories[path.Dir(file)]; !ok {
			directories[path.Dir(file)] = 0
		}
	}

	if verbose {
		logWithVerbose("Directories of changed files:")
		for directory := range directories {
			logfWithVerbose(" - %s", directory)
		}
	}

	goChanges := false
	otherChanges := false

	for directory := range directories {
		if directory == "." {
			otherChanges = true
			continue
		}

		if strings.Contains(directory, "cli/") {
			goChanges = true
			continue
		}
	}

	if goChanges && otherChanges {
		fmt.Println("both")
		os.Exit(0)
	}

	if goChanges {
		fmt.Println("go")
		os.Exit(0)
	}

	if otherChanges {
		fmt.Println("other")
		os.Exit(0)
	}
}

func logfWithVerbose(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}

func logWithVerbose(s string) {
	if verbose {
		log.Println(s)
	}
}
