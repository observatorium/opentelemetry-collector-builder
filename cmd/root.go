// Copyright 2020 OpenTelemetry Authors
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

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/observatorium/opentelemetry-collector-builder/internal/builder"
)

var cfgFile string
var cfg = builder.DefaultConfig()

// Execute is the main entrypoint for this application
func Execute() {
	cobra.OnInitialize(initConfig)

	cmd := &cobra.Command{
		Use:  "opentelemetry-collector-builder",
		Long: "OpenTelemetry Collector distribution builder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Validate(); err != nil {
				cfg.Logger.Error(err, "invalid configuration: %w", err)
				return nil
			}

			if err := cfg.ParseModules(); err != nil {
				cfg.Logger.Error(err, "invalid module configuration: %w", err)
				return nil
			}

			return builder.GenerateAndCompile(cfg)
		},
	}

	// the external config file
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.otelcol-builder.yaml)")

	// the distribution parameters, which we accept as CLI flags as well
	cmd.Flags().StringVar(&cfg.Distribution.ExeName, "name", "otelcol-custom", "The executable name for the OpenTelemetry Collector distribution")
	cmd.Flags().StringVar(&cfg.Distribution.LongName, "description", "Custom OpenTelemetry Collector distribution", "A descriptive name for the OpenTelemetry Collector distribution")
	cmd.Flags().StringVar(&cfg.Distribution.Version, "version", "1.0.0", "The version for the OpenTelemetry Collector distribution")
	cmd.Flags().StringVar(&cfg.Distribution.OtelColVersion, "otelcol-version", cfg.Distribution.OtelColVersion, "Which version of OpenTelemetry Collector to use as base")
	cmd.Flags().StringVar(&cfg.Distribution.OutputPath, "output-path", cfg.Distribution.OutputPath, "Where to write the resulting files")
	cmd.Flags().StringVar(&cfg.Distribution.Go, "go", "/usr/bin/go", "The Go binary to use during the compilation phase")
	cmd.Flags().StringVar(&cfg.Distribution.Module, "module", "github.com/jpkroehling/opentelemetry-collector-builder", "The Go module for the new distribution")

	// tie Viper to flags
	viper.BindPFlags(cmd.Flags())

	cmd.Execute()
}

func initConfig() {
	// a couple of Viper goodies, to make it easier to use env vars when flags are not desirable
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// load values from config file -- required for the modules configuration
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".otelcol-builder")
	}

	// load the config file
	if err := viper.ReadInConfig(); err == nil {
		cfg.Logger.Info("Using config file", "path", viper.ConfigFileUsed())
	}

	// convert Viper's internal state into our configuration object
	if err := viper.Unmarshal(&cfg); err != nil {
		cfg.Logger.Error(err, "failed to parse the config")
		return
	}
}
