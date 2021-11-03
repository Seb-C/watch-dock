package docker_compose

import (
	"github.com/Seb-C/watch-dock/pkg/context"
	"github.com/compose-spec/compose-go/types"
	yaml "gopkg.in/yaml.v2"
)

// Returns the built object of each docker-compose service that is built locally
// All the paths in this object are absolute
func GetServiceBuilds(ctx context.Context) (map[string]types.BuildConfig, error) {
	output, err := CallOnce(ctx, "config")
	if err != nil {
		return nil, err
	}

	config := struct{
		Services map[string]types.ServiceConfig `yaml:"services"`
	}{}
	if err = yaml.Unmarshal(output, &config); err != nil {
		return nil, err
	}

	builds := make(map[string]types.BuildConfig, 0)
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Build != nil {
			builds[serviceName] = *serviceConfig.Build
		}
	}

	return builds, err
}
