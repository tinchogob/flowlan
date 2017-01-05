package flowlan


func Run(tasks ...*Task) map[string]interface{} {

	var deps []string
	for _, t := range tasks {
		deps = append(deps, t.Name)
	}

	tasks = append(tasks, &Task{
		Deps: deps,
		Fx: func(args map[string]interface{}) (interface{}, error) {
			return args, nil
		},
		out: []chan interface{}{make(chan interface{})},
	})

	for _, t := range tasks {
		if len(t.Deps) > 0 {
			for _, dep := range t.Deps {
				for _, t2 := range tasks {
					if dep != "" && dep == t2.Name {
						pipe := make(chan interface{})
						t.in = append(t.in, pipe)
						t2.out = append(t2.out, pipe)
					}
				}
			}
		}
	}

	for _, iter := range tasks {
		t := iter
		go func() {
			args := make(map[string]interface{})
			if len(t.Deps) > 0 {
				for i, in := range t.in {
					dep := <-in
					args[t.Deps[i]] = dep
				}

			}

			v, err := t.Fx(args)

			if v != nil && err == nil && len(t.out) > 0 {
				for _, out := range t.out {
					out <- v
				}
			}
		}()
	}

	res := <-tasks[len(tasks)-1].out[0]

	return res.(map[string]interface{})
}

type Task struct {
	Name  string
	Deps  []string
	in    []chan interface{}
	out   []chan interface{}
	Fx    func(args map[string]interface{}) (interface{}, error)
}