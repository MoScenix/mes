package mtl

import (
	"os"

	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func TraceInit(serviceName, exportEndpoint string) provider.OtelProvider {
	restoreEndpoint := unsetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer restoreEndpoint()
	restoreTracesEndpoint := unsetEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	defer restoreTracesEndpoint()

	opts := []provider.Option{
		provider.WithServiceName(serviceName),
		provider.WithInsecure(),
		provider.WithEnableMetrics(false),
	}
	if exportEndpoint != "" {
		opts = append(opts, provider.WithExportEndpoint(exportEndpoint))
	}
	p := provider.NewOpenTelemetryProvider(opts...)
	return p
}

func unsetEnv(name string) func() {
	value, ok := os.LookupEnv(name)
	os.Unsetenv(name)
	return func() {
		if ok {
			os.Setenv(name, value)
		}
	}
}
