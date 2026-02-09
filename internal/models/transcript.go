package models

import "time"

// TaskTranscript is the per-task JSON file written to the transcript directory.
type TaskTranscript struct {
	TaskID      string                   `json:"task_id"`
	TaskName    string                   `json:"task_name"`
	Status      string                   `json:"status"`
	StartedAt   time.Time                `json:"started_at"`
	CompletedAt time.Time                `json:"completed_at"`
	DurationMs  int64                    `json:"duration_ms"`
	Prompt      string                   `json:"prompt"`
	FinalOutput string                   `json:"final_output"`
	Transcript  []TranscriptEntry        `json:"transcript"`
	Validations map[string]GraderResults `json:"validations,omitempty"`
	Session     SessionDigest            `json:"session"`
	ErrorMsg    string                   `json:"error_msg,omitempty"`
}
