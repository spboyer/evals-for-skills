package metrics

import "math"

// TriggerMetrics holds classification metrics for trigger accuracy.
type TriggerMetrics struct {
	TP        int     `json:"true_positives"`
	FP        int     `json:"false_positives"`
	TN        int     `json:"true_negatives"`
	FN        int     `json:"false_negatives"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1        float64 `json:"f1"`
	Accuracy  float64 `json:"accuracy"`
}

// TriggerResult pairs an expected trigger label with the actual outcome.
type TriggerResult struct {
	ShouldTrigger bool // expected: true = should activate
	DidTrigger    bool // actual: true = skill activated
}

// ComputeTriggerMetrics calculates precision, recall, F1, and accuracy
// from a set of trigger classification results.
// Returns nil when results is empty.
func ComputeTriggerMetrics(results []TriggerResult) *TriggerMetrics {
	if len(results) == 0 {
		return nil
	}

	var tp, fp, tn, fn int
	for _, r := range results {
		switch {
		case r.ShouldTrigger && r.DidTrigger:
			tp++
		case !r.ShouldTrigger && r.DidTrigger:
			fp++
		case !r.ShouldTrigger && !r.DidTrigger:
			tn++
		case r.ShouldTrigger && !r.DidTrigger:
			fn++
		}
	}

	total := tp + fp + tn + fn

	precision := safeDivide(float64(tp), float64(tp+fp))
	recall := safeDivide(float64(tp), float64(tp+fn))

	var f1 float64
	if precision+recall > 0 {
		f1 = 2 * precision * recall / (precision + recall)
	}

	accuracy := safeDivide(float64(tp+tn), float64(total))

	return &TriggerMetrics{
		TP:        tp,
		FP:        fp,
		TN:        tn,
		FN:        fn,
		Precision: roundTo4(precision),
		Recall:    roundTo4(recall),
		F1:        roundTo4(f1),
		Accuracy:  roundTo4(accuracy),
	}
}

func safeDivide(num, den float64) float64 {
	if den == 0 {
		return 0.0
	}
	return num / den
}

func roundTo4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
