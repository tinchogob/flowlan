package flowlan

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name  string
		tasks func() []*task
		results func() []int
		err   string
	}{
		//{
		//	name: "Task/Do/no action",
		//	tasks: func() []*task {
		//		return []*task{Task("Tincho")}
		//	},
		//},
		//{
		//	name: "Task/After/do parametersmismatch",
		//
		//},
		//{
		//	name: "Task/Do/no return value",
		//	tasks: func() []*task {
		//		return []*task{Task("Tincho").Do(func(){
		//			fmt.Println("hola")
		//		})}
		//	},
		//},
		{
			name: "Task/Do/1 return value",
			tasks: func() []*task {
				return []*task{Task("one").Do(func() int {
					return 1
				})}
			},
			results: func() []int {
				return []int{1}
			},
		},
		//{
		//	name: "Task/Do/multiple return values",
		//	tasks: func() []*task {
		//		return []*task{Task("multiple").Do(func() (int, int) {
		//			return 1, 2
		//		})}
		//	},
		//	results: func() []int {
		//		return []int{1, 2}
		//	},
		//},
		//{
		//	name: "Flow/After/wrong dependency",
		//	tasks: func() []*task {
		//		return []*task{
		//			Task("one").Do(func() int {
		//				return 1
		//			}),
		//			Task("two").After("chabon").Do(func(oneResult int) int {
		//				return 1+oneResult
		//			}),
		//		}
		//	},
		//	err: "error",
		//},
		//{
		//	name: "Flow/After/two many dependencies",
		//	tasks: func() []*task {
		//		return []*task{
		//			Task("one").Do(func() int {
		//				return 1
		//			}),
		//			Task("two").After("chabon").Do(func() int {
		//				return 1
		//			}),
		//		}
		//	},
		//	err: "error",
		//},
		//{
		//	name: "Flow/After/self dependency",
		//	tasks: func() []*task {
		//		return []*task{
		//			Task("one").Do(func() int {
		//				return 1
		//			}),
		//			Task("two").After("two").Do(func() int {
		//				return 1
		//			}),
		//		}
		//	},
		//	err: "error",
		//},
		//{
		//	name: "Flow/After/circular dependencies",
		//	tasks: func() []*task {
		//		return []*task{
		//			Task("one").After("three").Do(func() int {
		//				return 1
		//			}),
		//			Task("two").After("one").Do(func() int {
		//				return 1
		//			}),
		//			Task("three").After("two").Do(func() int {
		//				return 1
		//			}),
		//		}
		//	},
		//	err: "error",
		//},
		//{
		//	name: "Flow/After/deadlock dependencies",
		//	tasks: func() []*task {
		//		return []*task{
		//			Task("one").After("two").Do(func() int {
		//				return 1
		//			}),
		//			Task("two").After("one").Do(func() int {
		//				return 1
		//			}),
		//		}
		//	},
		//	err: "error",
		//},
		{
			name: "Flow/After/1 dependency",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() int {
						return 1
					}),
					Task("two").After("one").Do(func(oneResult int) int {
						return 1+oneResult
					}),
				}
			},
			results: func() []int {
				return []int{1, 2}
			},
		},
		{
			name: "Flow/After/2 dependencies",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() int {
						return 1
					}),
					Task("two").Do(func() int {
						return 1
					}),
					Task("three").After("one", "two").Do(func(oneResult, twoResult int) int {
						return oneResult+twoResult+1
					}),
				}
			},
			results: func() []int {
				return []int{1, 1, 3}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := Run(c.tasks()...)
			if err != nil {
				if c.err == err.Error() {
					t.Fatalf("%s: expected %s error but got %s", c.name, c.err, err.Error())
				}
			} else {
				if c.err != "" {
					t.Fatalf("%s: expected %s ", c.name, c.err)
				}

				results := c.results()

				if len(res) != len(results){
					t.Fatalf("%s: expected %s results but got %s", c.name, len(results), len(res))
				}

				for i, r := range res {
					if r != results[i] {
						t.Errorf("%s: expected %d but got %d", c.name, results[i], r)
					}
				}

				fmt.Println(res)
			}
		})
	}
}
