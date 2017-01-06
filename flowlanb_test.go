package flowlan

import (
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct{
		name string
		tasks func() []*task
		results func() []int
		err string
	}{
		{
			name: "Flow/one task",
			tasks: func()[]*task{
				return []*task{Task("1").Do(func(args map[string]interface{}) (interface{}, error){
					return 1,nil
				})}
			},
			results: func() []int {
				return []int{1}
			},
		},
		{
			name: "Flow/parallel tasks",
			tasks: func()[]*task{
				return []*task{
					Task("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1,nil
					}),
					Task("2").Do(func(args map[string]interface{}) (interface{}, error){
						return 2,nil
					}),
				}
			},
			results: func() []int {
				return []int{1, 2}
			},
		},
		{
			name: "Flow/series tasks",
			tasks: func()[]*task{
				return []*task{
					Task("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1, nil
					}),
					Task("2").After("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["1"].(int), nil
					}),
					Task("3").After("2").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["2"].(int), nil
					}),
				}
			},
			results: func() []int {
				return []int{1, 2, 3}
			},
		},
		{
			name: "Flow/many dependencies",
			tasks: func()[]*task{
				return []*task{
					Task("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1, nil
					}),
					Task("2").After("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["1"].(int), nil
					}),
					Task("3").After("1","2").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["1"].(int)+args["2"].(int), nil
					}),
				}
			},
			results: func() []int {
				return []int{1, 2, 4}
			},
		},
		{
			name: "Flow/zero value return",
			tasks: func()[]*task{
				return []*task{
					Task("1").Do(func(args map[string]interface{}) (interface{}, error){
						return nil, nil
					}),
					Task("2").After("1").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["1"].(int), nil
					}),
					Task("3").After("1","2").Do(func(args map[string]interface{}) (interface{}, error){
						return 1+args["1"].(int)+args["2"].(int), nil
					}),
				}
			},
			results: func() []int {
				return []int{1, 2, 4}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tasks := c.tasks()
			res, err := Run(tasks...)
			if err != nil {
				if err.Error() != c.err {
					t.Errorf("%s: expected %s but got %s", c.name, c.err, err.Error())
				}
			} else {
				if c.err != "" {
					t.Errorf("%s: expected %s", c.name, c.err)
				}

				for i, task := range tasks {
					if res[task.name].(int) != c.results()[i] {
						t.Errorf("%s: expected %d but got %d", c.name, c.results()[i], res[task.name].(int))
					}
				}

			}
		})
	}
}
