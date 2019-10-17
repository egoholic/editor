package creating

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

	Value struct {
	}
)
