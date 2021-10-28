package dockerfile

import (
	"bytes"
	"github.com/compose-spec/compose-go/types"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"os"
	"path/filepath"
	"strings"
)

// Returns all the local dependencies of the given dockerfile
// A dependency is a source params in a COPY or ADD instruction.
// The returned data is not transformed, it's a filepath.Match pattern
// whose meaning is relative to the context of the build.
func GetRawDependencies(dockerfilePath string) ([]string, error) {
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

// Transforms the given filepath.Match patterns in a list of existing files
// on the current filesystem. The output is still relative to the context.
func GetContextualizedFilePaths(context string, patterns []string) ([]string, error) {
	originalWorkDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err := os.Chdir(context); err != nil {
		return nil, err
	}
	defer os.Chdir(originalWorkDir)

	matchingPaths := make([]string, 0)
	for _, dependencyPath := range patterns {
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

func UncontextualizePaths(context string, paths []string) ([]string) {
	outPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		outPaths = append(outPaths, filepath.Join(context, path))
	}

	return outPaths
}

func GetAbsoluteDockerfileDependencies(build types.BuildConfig) ([]string, error) {
	rawDependencies, err := GetRawDependencies(build.Dockerfile)
	if err != nil {
		return nil, err
	}

	contextualizedPaths, err := GetContextualizedFilePaths(build.Context, rawDependencies)
	if err != nil {
		return nil, err
	}

	paths := UncontextualizePaths(build.Context, contextualizedPaths)

	return paths, nil
}
