package checks

// CheckStatus represents the three-tier status model used by score checks.
type CheckStatus string

const (
	// StatusOK indicates the check passes.
	StatusOK CheckStatus = "ok"
	// StatusOptimal indicates the check meets recommended best practice.
	StatusOptimal CheckStatus = "optimal"
	// StatusWarning indicates a potential issue was detected.
	StatusWarning CheckStatus = "warning"
)

// StatusHolder is implemented by check Data types that carry a CheckStatus.
type StatusHolder interface {
	GetStatus() CheckStatus
}
