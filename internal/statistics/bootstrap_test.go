package statistics

import (
	"math"
	"testing"
)

func TestBootstrapCI_EmptyScores(t *testing.T) {
	ci := BootstrapCI(nil, 0.95)
	if ci.Mean != 0.0 || ci.Lower != 0.0 || ci.Upper != 0.0 {
		t.Errorf("expected zero CI for empty input, got %+v", ci)
	}
	if ci.NumBootstraps != 0 {
		t.Errorf("expected 0 bootstraps for empty input, got %d", ci.NumBootstraps)
	}
}

func TestBootstrapCI_SingleValue(t *testing.T) {
	ci := BootstrapCI([]float64{0.75}, 0.95)
	if ci.Mean != 0.75 || ci.Lower != 0.75 || ci.Upper != 0.75 {
		t.Errorf("expected degenerate CI for single value, got %+v", ci)
	}
}

func TestBootstrapCI_IdenticalValues(t *testing.T) {
	ci := BootstrapCIWithSeed([]float64{0.5, 0.5, 0.5, 0.5}, 0.95, 42)
	if math.Abs(ci.Lower-0.5) > 1e-9 || math.Abs(ci.Upper-0.5) > 1e-9 {
		t.Errorf("expected CI [0.5, 0.5] for identical values, got [%f, %f]", ci.Lower, ci.Upper)
	}
}

func TestBootstrapCI_KnownDistribution(t *testing.T) {
	// 10 scores with known mean ~0.5
	scores := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}
	ci := BootstrapCIWithSeed(scores, 0.95, 42)

	if ci.Mean < 0.54 || ci.Mean > 0.56 {
		t.Errorf("expected mean ~0.55, got %f", ci.Mean)
	}
	if ci.Lower >= ci.Mean {
		t.Errorf("lower bound %f should be < mean %f", ci.Lower, ci.Mean)
	}
	if ci.Upper <= ci.Mean {
		t.Errorf("upper bound %f should be > mean %f", ci.Upper, ci.Mean)
	}
	if ci.Lower < 0 || ci.Upper > 1.0 {
		t.Errorf("CI should be within [0, 1] for these scores, got [%f, %f]", ci.Lower, ci.Upper)
	}
	if ci.NumBootstraps != DefaultBootstrapIterations {
		t.Errorf("expected %d bootstraps, got %d", DefaultBootstrapIterations, ci.NumBootstraps)
	}
	if ci.ConfidenceLevel != 0.95 {
		t.Errorf("expected confidence level 0.95, got %f", ci.ConfidenceLevel)
	}
}

func TestBootstrapCI_CIContainsMean(t *testing.T) {
	scores := []float64{0.3, 0.5, 0.7, 0.4, 0.6}
	ci := BootstrapCIWithSeed(scores, 0.95, 123)

	if ci.Lower > ci.Mean || ci.Upper < ci.Mean {
		t.Errorf("CI [%f, %f] should contain mean %f", ci.Lower, ci.Upper, ci.Mean)
	}
}

func TestBootstrapCI_NarrowerAtHigherN(t *testing.T) {
	small := []float64{0.3, 0.5, 0.7}
	large := []float64{0.3, 0.4, 0.5, 0.6, 0.7, 0.3, 0.4, 0.5, 0.6, 0.7,
		0.3, 0.4, 0.5, 0.6, 0.7, 0.3, 0.4, 0.5, 0.6, 0.7}

	ciSmall := BootstrapCIWithSeed(small, 0.95, 42)
	ciLarge := BootstrapCIWithSeed(large, 0.95, 42)

	widthSmall := ciSmall.Upper - ciSmall.Lower
	widthLarge := ciLarge.Upper - ciLarge.Lower

	if widthLarge >= widthSmall {
		t.Errorf("larger sample should yield narrower CI: small=%f, large=%f", widthSmall, widthLarge)
	}
}

func TestIsSignificant(t *testing.T) {
	tests := []struct {
		name string
		ci   ConfidenceInterval
		want bool
	}{
		{"both positive", ConfidenceInterval{Lower: 0.1, Upper: 0.5}, true},
		{"both negative", ConfidenceInterval{Lower: -0.5, Upper: -0.1}, true},
		{"crosses zero", ConfidenceInterval{Lower: -0.1, Upper: 0.3}, false},
		{"lower at zero", ConfidenceInterval{Lower: 0.0, Upper: 0.5}, false},
		{"upper at zero", ConfidenceInterval{Lower: -0.3, Upper: 0.0}, false},
		{"both zero", ConfidenceInterval{Lower: 0.0, Upper: 0.0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSignificant(tt.ci)
			if got != tt.want {
				t.Errorf("IsSignificant(%+v) = %v, want %v", tt.ci, got, tt.want)
			}
		})
	}
}

func TestNormalizedGain(t *testing.T) {
	tests := []struct {
		name      string
		pre, post float64
		want      float64
	}{
		{"basic gain", 0.4, 0.7, 0.5}, // (0.7-0.4)/(1-0.4) = 0.3/0.6 = 0.5
		{"no change", 0.5, 0.5, 0.0},
		{"full gain", 0.5, 1.0, 1.0},
		{"pre at ceiling", 1.0, 1.0, 0.0},
		{"low to high", 0.0, 0.5, 0.5},          // (0.5-0.0)/(1-0.0) = 0.5
		{"high pre small gain", 0.9, 0.95, 0.5}, // (0.95-0.9)/(1-0.9) = 0.05/0.1 = 0.5
		{"zero to one", 0.0, 1.0, 1.0},
		{"negative gain", 0.5, 0.3, -0.4}, // (0.3-0.5)/(1-0.5) = -0.2/0.5 = -0.4
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizedGain(tt.pre, tt.post)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("NormalizedGain(%f, %f) = %f, want %f", tt.pre, tt.post, got, tt.want)
			}
		})
	}
}

func TestBootstrapCI_Deterministic(t *testing.T) {
	scores := []float64{0.2, 0.4, 0.6, 0.8}
	ci1 := BootstrapCIWithSeed(scores, 0.95, 99)
	ci2 := BootstrapCIWithSeed(scores, 0.95, 99)

	if ci1.Lower != ci2.Lower || ci1.Upper != ci2.Upper {
		t.Errorf("same seed should produce identical CIs: %+v vs %+v", ci1, ci2)
	}
}

func TestBootstrapCI_DifferentConfidenceLevels(t *testing.T) {
	scores := []float64{0.1, 0.3, 0.5, 0.7, 0.9, 0.2, 0.4, 0.6, 0.8, 1.0}
	ci90 := BootstrapCIWithSeed(scores, 0.90, 42)
	ci99 := BootstrapCIWithSeed(scores, 0.99, 42)

	width90 := ci90.Upper - ci90.Lower
	width99 := ci99.Upper - ci99.Lower

	if width99 <= width90 {
		t.Errorf("99%% CI should be wider than 90%%: 90%%=%f, 99%%=%f", width90, width99)
	}
}
