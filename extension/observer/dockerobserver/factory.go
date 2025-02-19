// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dockerobserver // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/dockerobserver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	// The value of extension "type" in configuration.
	typeStr config.Type = "docker_observer"
)

// NewFactory should be called to create a factory with default values.
func NewFactory() component.ExtensionFactory {
	return component.NewExtensionFactory(
		typeStr,
		createDefaultConfig,
		createExtension,
		component.StabilityLevelBeta,
	)
}

func createDefaultConfig() config.Extension {
	return &Config{
		ExtensionSettings: config.NewExtensionSettings(config.NewComponentID(typeStr)),
		Endpoint:          "unix:///var/run/docker.sock",
		Timeout:           5 * time.Second,
		CacheSyncInterval: 60 * time.Minute,
		DockerAPIVersion:  defaultDockerAPIVersion,
	}
}

func createExtension(
	_ context.Context,
	settings component.ExtensionCreateSettings,
	cfg config.Extension,
) (component.Extension, error) {
	config := cfg.(*Config)
	return newObserver(settings.Logger, config)
}
