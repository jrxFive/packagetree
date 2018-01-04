package repomanager

const INDEX string = "INDEX"

type IndexOperator struct {
	instruction *instruction
	repo        *Repo
}

func NewIndexOperator(instruction *instruction, repo *Repo) *IndexOperator {
	return &IndexOperator{
		instruction: instruction,
		repo:        repo,
	}
}

func (o IndexOperator) Run() (string, error) {

	o.repo.mu.Lock()
	err := o.repo.backend.Add(o.instruction.packageName, o.instruction.packageDependencies...)
	o.repo.mu.Unlock()

	if err != nil {
		return FAIL, nil
	}
	return OK, nil

}

func (o IndexOperator) GetCommand() string {
	return o.instruction.cmd
}
