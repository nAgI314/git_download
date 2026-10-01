package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: git download <repository> [directory]")
		os.Exit(1)
	}

	repository := os.Args[1]

	directory := ""
	if len(os.Args) >= 3 {
		directory = os.Args[2]
	} else {
		name := strings.TrimRight(repository, "/")
		name = filepath.Base(name)
		name = strings.TrimSuffix(name, ".git")
		directory = name
	}

	fmt.Printf("Downloading %s...\n", repository)

	cmd := exec.Command("git", "clone", repository, directory)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	gitDir := filepath.Join(directory, ".git")

	if err := os.RemoveAll(gitDir); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove .git: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloaded to %s\n", directory)
}
