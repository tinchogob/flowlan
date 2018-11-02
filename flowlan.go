package flowlan

import (
	"context"
	"fmt"
	"reflect"
)

// Debug controls if flowland logs some useful debugging info to stdout
var Debug = true

func log(format string, a ...interface{}) {
	if Debug {
		fmt.Printf(format+"\n", a...)
	}
}

var nop interface{} = func() {}

// Task defines a task to be runned by flowlan
// Each Task will run on a separate go routine
type Task struct {
	Name string
	in   []*dependency
	out  []*dependency
	fx   fx
}

type dependency struct {
	name string
	res  chan interface{}
}

type fx func(map[string]interface{}) (interface{}, error)

// Run runs tasks as soon as each dependencies finish
func Run(tasks ...*Task) error {
	return run(context.Background(), tasks...)
}

// RunWithContext runs tasks as soon as each dependencies finish
func RunWithContext(ctx context.Context, tasks ...*Task) error {
	return run(ctx, tasks...)
}

func run(ctx context.Context, tasks ...*Task) error {
	plumb(tasks)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errors := make(chan error)
	var tasksPending []chan struct{}
	for _, task := range tasks {
		tasksPending = append(tasksPending, task.run(ctx, errors))
	}

	for _, pendingTask := range tasksPending {
		select {
		case <-ctx.Done():
		case <-pendingTask:
		case err := <-errors:
			cancel()
			return err
		}
	}
	return nil
}

// Defines the task name to be runned by flowlan
func Step(name string) *Task {
	return &Task{
		Name: name,
	}
}

func (t *Task) After(deps ...string) *Task {
	for _, dep := range deps {
		t.in = append(t.in, &dependency{dep, nil})
	}

	return t
}

func (t *Task) IDo(f fx) *Task {
	t.fx = f
	return t
}

func (t *Task) Do(fx interface{}) *Task {
	t.fx = t.toFx(reflect.ValueOf(fx))
	return t
}

func plumb(tasks []*Task) {
	for _, task := range tasks {
		for _, inDep := range task.in {
			for _, aTask := range tasks {
				if aTask.Name == inDep.name {
					log("connecting %s in with %s out", task.Name, inDep.name)
					pipe := make(chan interface{})
					inDep.res = pipe
					aTask.out = append(aTask.out, &dependency{task.Name, pipe})
				}
			}
		}
	}
}

func (t *Task) run(ctx context.Context, errors chan error) chan struct{} {
	done := make(chan struct{})
	go func() {
		args := make(map[string]interface{})

		for _, inDep := range t.in {
			log("%s is waiting for dependency %s", t.Name, inDep.name)
			select {
			case args[inDep.name] = <-inDep.res:
			case <-ctx.Done():
				return
			}
		}

		if t.fx == nil {
			close(done)
			return
		}

		res, err := t.fx(args)

		if err != nil {
			select {
			case <-ctx.Done():
			case errors <- err:
			}
			return
		}

		for _, outDep := range t.out {
			log("%s sending: %v to %s", t.Name, res, outDep.name)
			select {
			case <-ctx.Done():
			case outDep.res <- res:
				close(outDep.res)
			}
		}
		close(done)
	}()
	return done
}

func (t *Task) toFx(fx reflect.Value) fx {
	return func(arguments map[string]interface{}) (interface{}, error) {

		args := []reflect.Value{}

		for i, dep := range t.in {
			log("%s args are: %v", t.Name, arguments[dep.name].([]interface{}))
			argsAsArray := arguments[dep.name].([]interface{})
			for _, v := range argsAsArray {
				vDep := reflect.ValueOf(v)
				if vDep.IsValid() {
					args = append(args, vDep)
				} else {
					args = append(args, reflect.Zero(reflect.TypeOf(fx.Type().In(i))))
				}
			}
		}

		log("%s is calling its reflect fx with: %v", t.Name, args)
		ret := fx.Call(args)

		var res []interface{}
		var err error
		var ok bool

		for i, r := range ret {
			//last return value must be of type error
			if i == len(ret)-1 && r.IsValid() && r.Interface() != nil {
				err, ok = r.Interface().(error)
				if !ok {
					res = append(res, r.Interface())
					err = nil
				}
			} else if i < len(ret)-1 {
				res = append(res, r.Interface())
			}
		}

		return res, err
	}
}
