package creating

type (
	Task struct {
		Title              string
		Goal               string
		Description        string
		AcceptanceCriteria []string
	}

	Value struct {
		task *Task
	}

	TaskSaver interface {
		SaveTask(*Task) (string, error)
	}
)

func New() *Value {
	return &Value{}
}
