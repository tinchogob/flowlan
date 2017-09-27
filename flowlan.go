package flowlan

import (
	"context"
	"fmt"
	"reflect"
)

var debug bool = true

func log(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", a...)
	}
}

var nop interface{} = func() {}

type task struct {
	name string
	in   []*dependency
	out  []*dependency
	fxi  func(map[string]interface{}) (interface{}, error)
}

type dependency struct {
	name string
	res  chan interface{}
}

//Runs tasks as soon as all its dependencies are ready.
//Returns a map[string]interface{} with each task results
func Run(tasks ...*task) (map[string]interface{}, error) {
	plumb(tasks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resChs := make(map[string]chan interface{})
	errors := make(chan error)

	res := make(map[string]interface{})

	for _, task := range tasks {
		resChs[task.name] = task.irun(ctx, errors)
	}

	for _, task := range tasks {
		select {
		case res[task.name] = <-resChs[task.name]:
		case err := <-errors:
			cancel()
			return nil, err
		}
	}

	return res, nil
}

//The Task to run
func Task(name string) *task {
	return &task{
		name: name,
	}
}

func (t *task) After(deps ...string) *task {
	for _, dep := range deps {
		t.in = append(t.in, &dependency{dep, nil})
	}

	return t
}

func (t *task) IDo(fx func(map[string]interface{}) (interface{}, error)) *task {
	t.fxi = fx
	return t
}

func (t *task) Do(fx interface{}) *task {
	t.fxi = t.toFxi(reflect.ValueOf(fx))
	return t
}

func plumb(tasks []*task) {
	for _, task := range tasks {
		log("plumbing task: %s ", task.name)
		for _, inDep := range task.in {
			log("%s depends on %s", task.name, inDep.name)
			for _, aTask := range tasks {
				if aTask.name == inDep.name {
					log("connecting %s in with %s out", task.name, inDep.name)
					pipe := make(chan interface{})
					inDep.res = pipe
					aTask.out = append(aTask.out, &dependency{task.name, pipe})
				}
			}
		}
	}
}

func (t *task) irun(ctx context.Context, errors chan error) chan interface{} {
	resCh := make(chan interface{})
	go func() {
		defer close(resCh)

		args := make(map[string]interface{})

		for _, inDep := range t.in {
			log("%s is waiting for dependency %s", t.name, inDep.name)
			select {
			case args[inDep.name] = <-inDep.res:
			case <-ctx.Done():
				return
			}
		}

		if t.fxi == nil {
			return
		}

		res, err := t.fxi(args)

		if err != nil {
			select {
			case <-ctx.Done():
			case errors <- err:
			}
			return
		}

		for _, outDep := range t.out {
			log("%s sending: %v to", t.name, res, outDep.name)
			select {
			case <-ctx.Done():
			case outDep.res <- res:
				close(outDep.res)
				resCh <- res
			}
		}

		resCh <- res
	}()

	return resCh
}

func (t *task) toFxi(fx reflect.Value) func(map[string]interface{}) (interface{}, error) {
	return func(arguments map[string]interface{}) (interface{}, error) {

		//args := make([]reflect.Value, len(t.in))
		var args []reflect.Value

		for i, dep := range t.in {
			vDep := reflect.ValueOf(arguments[dep.name])
			if vDep.IsValid() {
				args = append(args, vDep)
			} else {
				args = append(args, reflect.Zero(reflect.TypeOf(fx.Type().In(i))))
			}
		}

		log("calling reflect fx with: %v", args)
		ret := fx.Call(args)

		var err error
		if ret[1].IsValid() && ret[1].Interface() != nil {
			err = ret[1].Interface().(error)
		}
		return ret[0].Interface(), err
	}
}
