package cmd

type Cmd struct {
	functionList []func (*Args, interface{}) (interface{}, error)
}

func (cmd *Cmd) Run() error {
	args := GetArgs()
	var chain interface{}
	var err error
	for _,f := range cmd.functionList {
		chain,err = f(args, chain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Cmd) AddStep(f func (*Args, interface{}) (interface{},error)) {
	cmd.functionList = append(cmd.functionList, f)
}

func NewCmd() *Cmd{
	return &Cmd{
	}
}