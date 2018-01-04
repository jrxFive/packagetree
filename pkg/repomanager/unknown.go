package repomanager

type UnknownOperator struct {
	instruction *instruction
}

func NewUnknownOperator() *UnknownOperator {
	return &UnknownOperator{
		instruction: &instruction{
			cmd: ERROR,
		},
	}
}

func (o UnknownOperator) Run() (string, error) {
	return ERROR, nil
}

func (o UnknownOperator) GetCommand() string {
	return o.instruction.cmd
}
