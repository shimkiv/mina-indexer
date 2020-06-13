package model

var (
	HeightStatusPending = "pending"  // Height is pending processing
	HeightStatusOK      = "success"  // Height has been processed succesfully
	HeightStatusError   = "error"    // Height has encountered an error
	HeightStatusSkip    = "skip"     // Height has been marked as skipped
	HeightStatusNoBlock = "no_block" // Height does not contain a block
)

type Height struct {
	Model

	Height     uint64  `json:"height"`
	Status     string  `json:"status"`
	RetryCount int     `json:"retry_count"`
	Error      *string `json:"error"`
}

type HeightStatusCount struct {
	Status string
	Num    int
}

// ShouldRetry returns true if height is retriable
func (h Height) ShouldRetry() bool {
	return h.Status == HeightStatusError && h.RetryCount < 3
}

// ShouldSkip returns true if height should be skipped
func (h Height) ShouldSkip() bool {
	return h.Status == HeightStatusSkip
}

// ResetForRetry clears the errors
func (h *Height) ResetForRetry() {
	h.Error = nil
	h.RetryCount++
}
