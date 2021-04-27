package joinpoint

type Worker interface {
	Run(f func())
}

type goWorker struct{}

func (g *goWorker) Run(f func()) {
	go f()
}

func NewGoWorker() Worker {
	return &goWorker{}
}
