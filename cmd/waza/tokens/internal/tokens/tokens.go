package tokens

import (
	"math"
)

// Counter counts tokens in text.
type Counter interface {
	Count(text string) int
}

// EstimatingCounter approximates token count as ~4 characters per token.
type EstimatingCounter struct {
	CharsPerToken int
}

func NewEstimatingCounter() *EstimatingCounter {
	return &EstimatingCounter{CharsPerToken: 4}
}

func (c *EstimatingCounter) Count(text string) int {
	if len(text) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(text)) / float64(c.CharsPerToken)))
}
