package main

import (
	"bytes"
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// TODO pass all the arguments to docker-compose
	// TODO parse the local docker-compose file to find the docker images and contexts

	dependenciesPaths, err := getLocalDependencies(os.Args[1]) // TODO find the docker file
	if err != nil {
		fmt.Printf("Error: %v", err) // TODO
		return
	}

	matchingPaths, err := getFilesMatching(".", dependenciesPaths) // TODO right root path
	if err != nil {
		fmt.Printf("Error: %v", err) // TODO
		return
	}

	fmt.Printf("Result: %#v\n", matchingPaths)

	// TODO watch the paths
	// TODO start the services
	// TODO rebuild and restart the services on change

	// TODO detect buildkit and enable it automatically if necessary
}

func getLocalDependencies(dockerfilePath string) ([]string, error) {
	dockerfile, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}

	ast, err := parser.Parse(bytes.NewBuffer(dockerfile))
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)

	// The tokenized dockerfile is a tree, but to get a list of instructions,
	// we only need to iterate the childrens of the root.
	for _, node := range ast.AST.Children {
		instruction, err := instructions.ParseInstruction(node)
		if err != nil {
			return nil, err
		}

		if copyInstruction, isCopy := instruction.(*instructions.CopyCommand); isCopy {
			if copyInstruction.From == "" {
				// COPY instruction from local files: adding it's source paths
				paths = append(paths, copyInstruction.SourcesAndDest.SourcePaths...)
			} else {
				// COPY instruction with a '--from': ignoring non-local files
				continue
			}
		} else if addInstruction, isAdd := instruction.(*instructions.AddCommand); isAdd {
			// ADD instruction: adding it's local (non-http) source paths
			for _, path := range addInstruction.SourcesAndDest.SourcePaths {
				if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
					paths = append(paths, path)
				} else {
					continue
				}
			}
		} else {
			// Neither COPY or ADD: Ignoring the line
			continue
		}
	}

	return paths, nil
}

func getFilesMatching(rootDir string, dependenciesPaths []string) ([]string, error) {
	matchingPaths := make([]string, 0)
	for _, dependencyPath := range dependenciesPaths {
		foundPaths, err := filepath.Glob(dependencyPath)
		if err != nil {
			return nil, err
		}

		if foundPaths != nil {
			matchingPaths = append(matchingPaths, foundPaths...)
		}
	}

	return matchingPaths, nil
}
