package flowlan

import (
	"context"
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
	fx   reflect.Value
}

type dependency struct {
	name string
	res  chan reflect.Value
}

// Task constructor
func Step(name string) *Task {
	return &Task{
		Name: name,
	}
}

// Defines the task dependencies. The task will run when all its dependencies are done
func (t *Task) After(deps ...string) *Task {
	for _, dep := range deps {
		t.in = append(t.in, &dependency{dep, nil})
	}

	return t
}

// The func to run
func (t *Task) Do(fx interface{}) *Task {
	t.fx = reflect.ValueOf(fx)
	return t
}

// Run runs tasks in order as soon as their dependencies
// finishes with its results as argunments
func Run(ctx context.Context, tasks ...*Task) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	err := plumb(tasks)
	if err != nil {
		return err
	}

	var tasksPending []chan struct{}
	for _, task := range tasks {
		if task.fx.IsValid() {
			log("running %s", task.Name)
			tasksPending = append(tasksPending, task.run(ctx))
		}
	}

	for _, pendingTask := range tasksPending {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pendingTask:
		}
	}

	return ctx.Err()
}

func plumb(tasks []*Task) error {
	for _, task := range tasks {

		//Skip nop tasks
		if !task.fx.IsValid() {
			continue
		}

		var args int
		for _, inDep := range task.in {

			//catch self dependencies
			if inDep.name == task.Name {
				return fmt.Errorf("invalid task definition: circular dependency on %s", task.Name)
			}

			var found bool
			for _, aTask := range tasks {
				if aTask.Name == inDep.name {
					found = true
					args += aTask.fx.Type().NumOut()
					log("connecting %s in with %s out", task.Name, inDep.name)
					inDep.res = make(chan reflect.Value)
					aTask.out = append(aTask.out, &dependency{task.Name, inDep.res})
				}
			}
			if !found {
				return fmt.Errorf("invalid task definition: %s unknown dependency %s", task.Name, inDep.name)
			}
		}

		if args != task.fx.Type().NumIn() {
			return fmt.Errorf("invalid task definition: %s has %d arguments but got only %d", task.Name, task.fx.Type().NumIn(), args)
		}
	}
	return nil
}

func (t *Task) run(ctx context.Context) chan struct{} {
	done := make(chan struct{})

	go func() {

		numIn := t.fx.Type().NumIn()
		args := []reflect.Value{}

		for depIndex, inDep := range t.in {
			argCount := 0
			log("%d: %s is waiting for dependency %s", depIndex, t.Name, inDep.name)
			for anInDepRes := range inDep.res {
				log("%s received arg number %d with value %v from dependecy %s", t.Name, argCount, anInDepRes, inDep.name)
				if argCount >= numIn {
					log("%s is skipping arg number %d from dependency %s", t.Name, argCount, inDep.name)
					continue
				} else {
					log("%s is appending arg number %d with %v from dependency %s", t.Name, argCount, anInDepRes, inDep.name)
					args = append(args, anInDepRes)
				}
				argCount++
			}
		}

		if ctx.Err() != nil {
			return
		}

		log("calling fx with %d/%d", len(args), numIn)
		res := t.fx.Call(args)

		if ctx.Err() != nil {
			return
		}

		for _, outDep := range t.out {
			log("%s sending: %v to %s", t.Name, res, outDep.name)
			for _, outDepRes := range res {
				outDep.res <- outDepRes
			}
			close(outDep.res)
		}

		close(done)
	}()

	return done
}
