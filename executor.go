package juck

type Runnable interface {
	Run()
}

type Executor interface {
	Execute(r Runnable)
}
