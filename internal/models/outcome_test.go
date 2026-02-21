package models

import (
	"math"
	"testing"
)

func TestComputeStdDev(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{name: "empty", values: []float64{}, want: 0.0},
		{name: "single value", values: []float64{0.5}, want: 0.0},
		{name: "identical values", values: []float64{0.8, 0.8, 0.8}, want: 0.0},
		{name: "known values", values: []float64{2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0}, want: 2.0},
		{name: "two values", values: []float64{0.0, 1.0}, want: 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeStdDev(tt.values)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("ComputeStdDev(%v) = %f, want %f", tt.values, got, tt.want)
			}
		})
	}
}

func TestComputeRunScore(t *testing.T) {
	tests := []struct {
		name string
		run  RunResult
		want float64
	}{
		{name: "no validations", run: RunResult{}, want: 0.0},
		{
			name: "single validation",
			run:  RunResult{Validations: map[string]GraderResults{"check": {Score: 0.75, Passed: true}}},
			want: 0.75,
		},
		{
			name: "multiple validations",
			run: RunResult{Validations: map[string]GraderResults{
				"a": {Score: 1.0, Passed: true},
				"b": {Score: 0.5, Passed: false},
			}},
			want: 0.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.run.ComputeRunScore()
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("ComputeRunScore() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestComputeWeightedRunScore(t *testing.T) {
	tests := []struct {
		name string
		run  RunResult
		want float64
	}{
		{name: "no validations", run: RunResult{}, want: 0.0},
		{
			name: "single validation default weight",
			run:  RunResult{Validations: map[string]GraderResults{"check": {Score: 0.75, Weight: 1.0}}},
			want: 0.75,
		},
		{
			name: "equal weights same as unweighted",
			run: RunResult{Validations: map[string]GraderResults{
				"a": {Score: 1.0, Weight: 1.0},
				"b": {Score: 0.5, Weight: 1.0},
			}},
			want: 0.75,
		},
		{
			name: "weighted favoring higher scorer",
			run: RunResult{Validations: map[string]GraderResults{
				"a": {Score: 1.0, Weight: 3.0},
				"b": {Score: 0.0, Weight: 1.0},
			}},
			want: 0.75, // (1.0*3 + 0.0*1) / (3+1) = 0.75
		},
		{
			name: "zero weight defaults to 1.0",
			run: RunResult{Validations: map[string]GraderResults{
				"a": {Score: 1.0, Weight: 0.0},
				"b": {Score: 0.5, Weight: 0.0},
			}},
			want: 0.75, // treated as equal weight 1.0
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.run.ComputeWeightedRunScore()
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("ComputeWeightedRunScore() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestAllValidationsPassed(t *testing.T) {
	tests := []struct {
		name string
		run  RunResult
		want bool
	}{
		{name: "no validations passes", run: RunResult{}, want: true},
		{
			name: "all passed",
			run:  RunResult{Validations: map[string]GraderResults{"a": {Passed: true}, "b": {Passed: true}}},
			want: true,
		},
		{
			name: "one failed",
			run:  RunResult{Validations: map[string]GraderResults{"a": {Passed: true}, "b": {Passed: false}}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.run.AllValidationsPassed()
			if got != tt.want {
				t.Errorf("AllValidationsPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestStatsStdDevScore(t *testing.T) {
	scores := []float64{1.0, 0.5, 0.5}
	got := ComputeStdDev(scores)

	mean := (1.0 + 0.5 + 0.5) / 3.0
	variance := ((1.0-mean)*(1.0-mean) + (0.5-mean)*(0.5-mean) + (0.5-mean)*(0.5-mean)) / 3.0
	want := math.Sqrt(variance)

	if math.Abs(got-want) > 1e-9 {
		t.Errorf("ComputeStdDev() = %f, want %f", got, want)
	}
}
