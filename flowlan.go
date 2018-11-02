package flowlan

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Debug controls if flowland logs some useful debugging info to stdout
var Debug = false

func log(format string, a ...interface{}) {
	if Debug {
		fmt.Printf(format+"\n", a...)
	}
}

// Task defines a task to be runned by flowlan
// Each Task will run on a separate go routine
type Task struct {
	Name string
	in   []*dependency
	out  []*dependency
	fx   interface{}
}

type dependency struct {
	name string
	res  chan reflect.Value
}

// Run runs tasks as soon as each dependencies finish
func Run(ctx context.Context, tasks ...*Task) error {
	err := plumb(tasks)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errors := make(chan error)
	var tasksPending []chan struct{}
	for _, task := range tasks {
		if task.fx != nil {
			tasksPending = append(tasksPending, task.run(ctx, errors))
		}
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

func (t *Task) Do(fx interface{}) *Task {
	t.fx = fx
	return t
}

func plumb(tasks []*Task) error {
	for _, task := range tasks {
		for _, inDep := range task.in {
			var found bool
			for _, aTask := range tasks {
				if aTask.Name == inDep.name {
					found = true
					log("connecting %s in with %s out", task.Name, inDep.name)
					inDep.res = make(chan reflect.Value)
					aTask.out = append(aTask.out, &dependency{task.Name, inDep.res})
				}
			}
			if !found || inDep.name == task.Name {
				return errors.New("invalid task definition")
			}
		}
	}
	return nil
}

var nilErrorValue = reflect.ValueOf(new(error)).Elem()

func (t *Task) run(ctx context.Context, errors chan error) chan struct{} {

	done := make(chan struct{})

	go func() {

		fx := reflect.ValueOf(t.fx)
		numIn := fx.Type().NumIn()

		args := []reflect.Value{}
		for depIndex, inDep := range t.in {
			argCount := 0
			log("%d: %s is waiting for dependency %s", depIndex, t.Name, inDep.name)
			for anInDepRes := range inDep.res {
				log("%s received arg number %d with value %v from dependecy %s", t.Name, argCount, anInDepRes, inDep.name)
				if argCount >= numIn {
					log("%s is skipping arg number %d from dependency %s", t.Name, argCount, inDep.name)
					continue
					// } else if fx.Type().In(i) == nilErrorValue.Type() {
					// 	log("%s is appending arg number %d with nilErrorValue from dependency %s", t.Name, i, inDep.name)
					// 	args = append(args, nilErrorValue)
				} else {
					log("%s is appending arg number %d with %v from dependency %s", t.Name, argCount, anInDepRes, inDep.name)
					args = append(args, anInDepRes)
				}
				argCount++
			}
		}

		log("calling fx with %d/%d", len(args), numIn)
		res := fx.Call(args)

		for _, outDep := range t.out {
			log("%s sending: %v to %s", t.Name, res, outDep.name)
			for _, outDepRes := range res {
				select {
				case <-ctx.Done():
				case outDep.res <- outDepRes:
				}
			}
			close(outDep.res)
		}

		close(done)
	}()
	return done
}
