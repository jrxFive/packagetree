package repomanager

const QUERY string = "QUERY"

type QueryOperator struct {
	instruction *instruction
	repo        *Repo
}

func NewQueryOperator(instruction *instruction, repo *Repo) *QueryOperator {
	return &QueryOperator{
		instruction: instruction,
		repo:        repo,
	}
}

func (o QueryOperator) Run() (string, error) {

	o.repo.mu.RLock()
	exists := o.repo.backend.Exists(o.instruction.packageName)
	o.repo.mu.RUnlock()

	if exists {
		return OK, nil
	}
	return FAIL, nil

}

func (o QueryOperator) GetCommand() string {
	return o.instruction.cmd
}
