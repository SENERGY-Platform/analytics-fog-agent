module github.com/SENERGY-Platform/analytics-fog-agent

//replace github.com/SENERGY-Platform/analytics-fog-lib => ../analytics-fog-lib

require (
	github.com/SENERGY-Platform/analytics-fog-lib v1.1.16
	github.com/SENERGY-Platform/go-service-base v0.13.0
	github.com/SENERGY-Platform/mgw-module-manager/aux-client v0.4.0
	github.com/SENERGY-Platform/mgw-module-manager/lib v0.4.0
	github.com/docker/distribution v2.8.3+incompatible
	github.com/docker/docker v26.1.1+incompatible
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/joho/godotenv v1.5.1
	github.com/y-du/go-log-level v1.0.0
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/SENERGY-Platform/go-base-http-client v0.0.2 // indirect
	github.com/SENERGY-Platform/go-service-base/job-hdl/lib v0.1.0 // indirect
	github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib v0.0.2 // indirect
	github.com/SENERGY-Platform/mgw-module-lib v0.19.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/y-du/go-env-loader v0.5.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.51.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)

go 1.22

toolchain go1.22.2
