package flowlan

import (
	"fmt"
	"reflect"
)

var debug bool = false

func log(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", a...)
	}
}

func plumb(tasks []*task) {
	for _, task := range tasks {
		if len(task.dependencies) > 0 {
			for _, dependencyName := range task.dependencies {
				for _, dependency := range tasks {
					if dependencyName == dependency.name {
						pipe := make(chan interface{})
						dependency.out = append(dependency.out, pipe)
						task.in = append(task.in, pipe)
						log("connecting %s out with %s in", dependency.name, task.name)
					}
				}
			}
		}
	}
}

func collector(tasks []*task) []chan interface{} {
	resCh := make([]chan interface{}, len(tasks))

	for i, task := range tasks {
		resCh[i] = make(chan interface{})
		task.out = append(task.out, resCh[i])
		log("setting result channel for %s out", task.name)
	}

	return resCh
}

func Run(tasks ...*task) ([]interface{}, error) {
	plumb(tasks)

	resCh := collector(tasks)

	for _, task := range tasks {
		go task.run()
	}

	res := make([]interface{}, 0)

	for _, rCh := range resCh {
		for r := range rCh{
			res = append(res, r)
		}
	}

	return res, nil
}

var nop interface{} = func(){}

func Task(name string) *task {
	t := &task{
		name: name,
	}
	return t.Do(nop)
}

func (t *task) run() {
	var args []reflect.Value
	log("%s dependencies are: %v", t.name, t.dependencies)
	for i, depName := range t.dependencies {
		log("%s is waiting for dependency %s", t.name, depName)
		for dep := range t.in[i] {
			log("%s received %v from dependency %s", t.name, dep, depName)

			vDep := reflect.ValueOf(dep)
			if vDep.IsValid() {
				args = append(args, vDep)
			} else {
				args = append(args, reflect.Zero(t.fx.Type().In(i)))
			}
		}
	}

	log("exec %s with args: %v", t.name, args)

	res := t.fx.Call(args)

	log("ended %s with res: %v", t.name, res)

	for _, out := range t.out {
		for i, r := range res {
			//technically no fx may return an invalid value as per stated in IsValid() docs
			if r.IsValid() {
				log("%s sending resInterface: %v", t.name, r.Interface())
				out <- r.Interface()
			} else {
				log("%s sending zeroType: %v", t.name, reflect.Zero(t.fx.Type().Out(i)))
				out <- reflect.Zero(t.fx.Type().Out(i))
			}
		}
		close(out)
	}
}

func (t *task) After(deps ...string) *task {
	t.dependencies = deps
	return t
}

func (t *task) Do(fx interface{}) *task {
	f := reflect.ValueOf(fx)
	t.fx = f

	return t
}

type task struct {
	name         string
	dependencies []string
	in           []chan interface{}
	out          []chan interface{}
	fx           reflect.Value
}
