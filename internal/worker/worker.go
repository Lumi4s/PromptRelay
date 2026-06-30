package worker

import "log"

type worker struct {
	id int8
}

func New() *worker {
	log.Println("Worker was successfully created")
	return &worker{
		id: 1,
	}
}

func (w *worker) Run() {

}
