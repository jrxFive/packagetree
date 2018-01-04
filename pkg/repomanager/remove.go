package repomanager

const REMOVE string = "REMOVE"

type RemoveOperator struct {
	instruction *instruction
	repo        *Repo
}

func NewRemoveOperator(instruction *instruction, repo *Repo) *RemoveOperator {
	return &RemoveOperator{
		instruction: instruction,
		repo:        repo,
	}
}

func (o RemoveOperator) Run() (string, error) {

	o.repo.mu.Lock()
	err := o.repo.backend.Remove(o.instruction.packageName)
	o.repo.mu.Unlock()

	if err != nil {
		return FAIL, nil
	}

	return OK, nil
}

func (o RemoveOperator) GetCommand() string {
	return o.instruction.cmd
}
