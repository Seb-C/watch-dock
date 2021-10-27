package dockerfile

import (
	"bytes"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"os"
	"strings"
)

func GetLocalDependencies(dockerfilePath string) ([]string, error) {
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
