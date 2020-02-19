package CCron

type Cmd interface {
	Run()
}

type CmdFunc func()

func (this CmdFunc) Run()  {
	 this()
}
