package main

import (
	"fmt"
	"path/filepath"
	"github.com/Seb-C/watch-dock/pkg/dockerfile"
	dockerCompose "github.com/Seb-C/watch-dock/pkg/docker_compose"
)

func main() {
	serviceBuilds, err := dockerCompose.GetServiceBuilds() // TODO
	if err != nil {
		fmt.Printf("Error: %v", err) // TODO
		return
	}
	fmt.Printf("Config: %v", out) // TODO

	// TODO check docker-compose version
	// TODO pass all the arguments to docker-compose
	// TODO parse the local docker-compose file to find the docker images and contexts
	// TODO unit tests

	dependenciesPaths, err := dockerfile.GetLocalDependencies("./build/docker/hello-world.Dockerfile") // TODO find the docker file
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

// TODO how to make this work with fsnotify?
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
