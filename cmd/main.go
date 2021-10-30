package main

import (
	"fmt"
	"github.com/Seb-C/watch-dock/pkg/context"
	"github.com/Seb-C/watch-dock/pkg/dockerfile"
	dockerCompose "github.com/Seb-C/watch-dock/pkg/docker_compose"
)

func main() {
	ctx := context.ParseContext()

	serviceBuilds, err := dockerCompose.GetServiceBuilds(ctx)
	if err != nil {
		fmt.Printf("Error: %v", err) // TODO
		return
	}

	// TODO test the arguments passed to docker-compose + convenient way to test
	// TODO check docker-compose version
	// TODO unit tests
	// TODO document and comment
	// TODO logging

	pathsToWatch := map[string]struct{}{}
	for _, serviceBuild := range serviceBuilds {
		dependencies, err := dockerfile.GetAbsoluteDockerfileDependencies(serviceBuild)
		if err != nil {
			fmt.Printf("Error: %v", err) // TODO
			return
		}

		for _, dependency := range dependencies {
			// Making all the paths unique
			pathsToWatch[dependency] = struct{}{}
		}
	}

	fmt.Printf("Result: %#v\n", pathsToWatch)

	// TODO watch the paths (how to make this work with fsnotify? is it recursive?)
	// TODO start the services
	// TODO rebuild and restart the services on change
	// TODO detect buildkit and enable it automatically if necessary
}
