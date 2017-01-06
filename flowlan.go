package flowlan

import (
	"fmt"
)

var debug bool = true

func log(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", a...)
	}
}

func Run(tasks ...*task) (map[string]interface{}, error) {

	var deps []string
	for _, t := range tasks {
		deps = append(deps, t.name)
	}

	//Append a final task that depends on all tasks and collects results
	tasks = append(tasks, &task{
		dependencies: deps,
		fx: func(args map[string]interface{}) (interface{}, error) {
			return args, nil
		},
		out: []chan interface{}{make(chan interface{})},
	})

	//do plumbing to connect dependant tasks in/out via channels
	for _, task := range tasks {
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

	for _, task := range tasks {
		go task.run()
	}

	res := <-tasks[len(tasks)-1].out[0]

	return res.(map[string]interface{}), nil
}

func Task(name string) *task {
	return &task{
		name: name,
	}
}

func (t *task) run() {
	args := make(map[string]interface{})
	for i, depName := range t.dependencies {
		log("%s waiting for dep %s", t.name, depName)
		dep := <-t.in[i]
		args[depName] = dep
	}

	log("exec %s", t.name)

	v, err := t.fx(args)

	log("ended %s", t.name)

	if v != nil && err == nil && len(t.out) > 0 {
		for _, out := range t.out {
			out <- v
			log("%s sent result", t.name)
		}
	}
}

func (t *task) After(deps ...string) *task {
	t.dependencies = deps
	return t
}

func (t *task) Do(f func(args map[string]interface{}) (interface{}, error)) *task {
	t.fx = f
	return t
}

type task struct {
	name         string
	dependencies []string
	in           []chan interface{}
	out          []chan interface{}
	fx           func(args map[string]interface{}) (interface{}, error)
}
