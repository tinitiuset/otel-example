package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/tinitiuset/otel-example/pkg/generator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	collectorEndpoint := getEnvOrDefault("OTEL_COLLECTOR_ENDPOINT", "localhost:4317")
	log.Printf("Using OpenTelemetry Collector at: %s", collectorEndpoint)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("wave-generator"),
			semconv.ServiceVersion("1.0.0"),
			semconv.ServiceInstanceID(hostname),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.DialContext(ctx, collectorEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatal(err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(15*time.Second),
			),
		),
	)
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter("wave-generator")

	waves := []string{
		"wave_1",
		"wave_2",
		"wave_3",
		"wave_4",
		"wave_5",
		"wave_6",
		"wave_7",
		"wave_8",
	}

	for _, name := range waves {
		g, err := generator.New(name, meter)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Started metric generator: %s (range: %.1f to %.1f, starting at %.1f, direction: %v)",
			name, g.GetMinValue(), g.GetMaxValue(), g.GetValue(), g.IsUp())
	}

	log.Printf("Metric generators started on instance %s. Updating every 15 seconds...", hostname)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nShutting down...")
}
