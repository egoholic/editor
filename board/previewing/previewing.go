package previewing

import "time"

type (
	Task struct {
		Title              string
		Context            string
		Description        string
		Goal               string
		AcceptanceCriteria []string
	}
	Draft struct {
		Content    string
		RawContent string
		Task
	}
	Approvement struct {
		ApprovedAt time.Time
		ApprovedBy string
		Comment    string
	}
	ApprovedDraft struct {
		Draft
		Approvement
	}

	Value struct {
	}

	TasksProvider interface {
		Tasks() ([]*Task, error)
	}

	DraftsProvider interface {
		Drafts() ([]*Draft, error)
	}

	ApprovedDraftsProvider interface {
		ApprovedDrafts() ([]*ApprovedDraft, error)
	}
)
