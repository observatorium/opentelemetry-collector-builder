# OpenTelemetry Collector builder

This program generates a custom OpenTelemetry Collector distribution based on a given configuration.

## Installation

```console
$ go install github.com/jpkroehling/opentelemetry-collector-builder
```

## Running

```console
$ otelcol-builder --config config.yaml
```

Use `otelcol-builder --help` to learn about which flags are available.

## Configuration

The configuration file is composed of two main sections: `dist` and `exporters`. All `dist` options can be specified via command line flags or environment variables: 

```console
$ otelcol-builder --name="my-otelcol"
```

```yaml
dist:
    module: otelcol-custom # the module name for the new distribution, following Go mod conventions. Optional, but recommended.
    name: otelcol-custom # the binary name. Optional.
    description: "Custom OpenTelemetry Collector distribution" # a long name for the application. Optional.
    output_path: /tmp/otelcol-distributionNNN # the path to write the output (sources and binary). Optional.
    version: "1.0.0" # the version for your custom OpenTelemetry Collector. Optional.
    go: "/usr/bin/go" # which Go binary to use to compile the generated sources. Optional.
exporters:
  - gomod: "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.9.0" # the Go module for the component. Required.
    import: "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter" # the import path for the component. Optional.
    name: "alibabacloudlogserviceexporter" # package name to use in the generated sources. Optional.
```
