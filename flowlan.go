package flowlan

import (
	"reflect"
	"fmt"
)

/*
TODO list
- Do: support for variadic functions
- Do: support for multiple return values
- Do: support for timeout
- General: error Handling

*/

var debug bool = true

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

func collector(tasks []*task) []chan interface {} {
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

	res := make([]interface{}, len(tasks))

	for i, rCh := range resCh {
		res[i] = <- rCh
	}

	return res, nil
}

func Task(name string) *task {
	return &task{
		name: name,
	}
}

func(t *task) run(){
	var args []reflect.Value
	if len(t.dependencies) > 0 {
		for i, depName := range t.dependencies {
			log("%s is waiting for dependency %s", t.name, depName)
			dep := <-t.in[i]
			log("%s received %v from dependency %s", t.name, dep, depName)

			inZeroType := reflect.Zero(t.fx.Type().In(i))
			if dep != inZeroType {
				args = append(args, reflect.ValueOf(dep))
			} else {
				args = append(args, inZeroType)
			}
		}

	}

	log("exec %s", t.name)
	res := t.fx.Call(args)
	log("ended %s", t.name)

	for i, r := range res {
		outZeroType := reflect.Zero(t.fx.Type().Out(i))
		if i == 0 && r != outZeroType {
			for _, out := range t.out {
				out <- res[0].Interface()
			}
		}
	}
}

func (t *task)After(deps ...string) *task {
	for _, dep := range deps {
		if dep == "" {
			panic("wrong task definition")
		}
	}

	t.dependencies = deps
	return t
}

func (t *task) Do(fx interface{}) *task {
	f := reflect.ValueOf(fx)

	if f.Type().NumIn() != len(t.dependencies) {
		panic(fmt.Sprintf("wrong task definition: %s: has %d dependencies but %d args\n", t.name, len(t.dependencies), f.Type().NumIn()))
	}

	t.fx = f

	return t
}

type task struct {
	name  string
	dependencies  []string
	in    []chan interface{}
	out   []chan interface{}
	fx    reflect.Value
}