package generator

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Generator struct {
	name     string
	minValue float64
	maxValue float64
	value    float64
	isUp     bool
	meter    metric.Meter
	wave     metric.Float64ObservableGauge
}

func New(name string, meter metric.Meter) (*Generator, error) {
	min := rand.Float64()*200 - 100
	max := rand.Float64()*200 - 100

	if min > max {
		min, max = max, min
	}

	startValue := min + rand.Float64()*(max-min)

	g := &Generator{
		name:     name,
		minValue: min,
		maxValue: max,
		value:    startValue,
		isUp:     rand.Float64() > 0.5,
		meter:    meter,
	}

	var err error
	g.wave, err = meter.Float64ObservableGauge(
		name,
		metric.WithDescription("Bouncing value metric"),
	)
	if err != nil {
		return nil, err
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			g.observe(o)
			return nil
		},
		g.wave,
	)
	if err != nil {
		return nil, err
	}

	go g.run()
	return g, nil
}

func (g *Generator) GetMinValue() float64 {
	return g.minValue
}

func (g *Generator) GetMaxValue() float64 {
	return g.maxValue
}

func (g *Generator) GetValue() float64 {
	return g.value
}

func (g *Generator) IsUp() bool {
	return g.isUp
}

func (g *Generator) observe(o metric.Observer) {
	attrs := []attribute.KeyValue{
		attribute.Float64("min_value", g.minValue),
		attribute.Float64("max_value", g.maxValue),
	}

	o.ObserveFloat64(g.wave, g.value, metric.WithAttributes(attrs...))
}

func (g *Generator) run() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if g.isUp {
			g.value++
			if g.value >= g.maxValue {
				g.isUp = false
			}
		} else {
			g.value--
			if g.value <= g.minValue {
				g.isUp = true
			}
		}
	}
}
