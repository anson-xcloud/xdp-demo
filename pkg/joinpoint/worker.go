package joinpoint

type Worker interface {
	Run(f func())
}

type goWorker struct {
	recover func(interface{})
}

func (g *goWorker) Run(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if g.recover != nil {
					g.recover(r)
				}
			}
		}()

		f()
	}()
}

func NewGoWorker(recover func(interface{})) Worker {
	return &goWorker{recover: recover}
}

type syncWorker struct{}

func (w *syncWorker) Run(f func()) {
	f()
}

func NewSyncWorker() Worker {
	return &syncWorker{}
}
