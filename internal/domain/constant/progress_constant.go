package constant

type ProgressStatus string

const (
    StatusNotStarted ProgressStatus = "not_started"
    StatusInProgress ProgressStatus = "in_progress"
    StatusCompleted  ProgressStatus = "completed"
    StatusFailed     ProgressStatus = "failed"
)
