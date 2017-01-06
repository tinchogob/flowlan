package flowlan

import (
	"io"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	mock := new(task)

	cases := []struct {
		name    string
		tasks   func() []*task
		results func() []interface{}
		err     string
	}{
		{
			name: "Task/Do/no action",
			tasks: func() []*task {
				return []*task{Task("Tincho")}
			},
			results: func() []interface{} {
				return []interface{}{}
			},
		},
		{
			name: "Task/Do/no return value",
			tasks: func() []*task {
				return []*task{Task("Tincho").Do(func() {
					//Do some work
				})}
			},
			results: func() []interface{} {
				return []interface{}{}
			},
		},
		{
			name: "Task/Do/1 return value",
			tasks: func() []*task {
				return []*task{Task("one").Do(func() int {
					return 1
				})}
			},
			results: func() []interface{} {
				return []interface{}{1}
			},
		},
		{
			name: "Task/Do/zero value return",
			tasks: func() []*task {
				return []*task{Task("Tincho").Do(func() int {
					return 0
				})}
			},
			results: func() []interface{} {
				return []interface{}{0}
			},
		},
		{
			name: "Task/Do/zero value return interface value",
			tasks: func() []*task {
				return []*task{Task("Tincho").Do(func() interface{} {
					return nil
				})}
			},
			results: func() []interface{} {
				return []interface{}{nil}
			},
		},
		{
			name: "Task/Do/multiple return values",
			tasks: func() []*task {
				return []*task{Task("multiple").Do(func() (int, int) {
					return 1, 2
				})}
			},
			results: func() []interface{} {
				return []interface{}{1, 2}
			},
		},
		{
			name: "Task/Do/multiple dependency return values",
			tasks: func() []*task {
				return []*task{
					Task("multiple_1").Do(func() (int, int) {
						return 1, 1
					}),
					Task("dependant").After("multiple_1").Do(func(d1, d2 int) int {
						return d1 + d2 + 1
					}),
				}
			},
			results: func() []interface{} {
				return []interface{}{1, 1, 3}
			},
		},
		{
			name: "Task/Do/error return value",
			tasks: func() []*task {
				return []*task{
					Task("multiple_1").Do(func() error {
						return io.ErrNoProgress
					}),
					Task("multiple_2").After("multiple_1").Do(func(e1 error) int {
						return 1
					}),
				}
			},
			results: func() []interface{} {
				return []interface{}{io.ErrNoProgress, 1}
			},
		},
		{
			name: "Task/Do/zero return value error",
			tasks: func() []*task {
				return []*task{
					Task("zre_1").Do(func() error {
						return nil
					}),
					Task("zre_2").After("zre_1").Do(func(e1 error) int {
						return 1
					}),
				}
			},
			results: func() []interface{} {
				return []interface{}{nil, 1}
			},
		},
		{
			name: "Task/Do/return value ptr",
			tasks: func() []*task {
				return []*task{
					Task("s_1").Do(func() *task {
						mock.name = "tinchogob"
						return mock
					}),
					Task("s_2").After("s_1").Do(func(s1 *task) string {
						return "el nombre es: " + s1.name
					}),
				}
			},
			results: func() []interface{} {
				return []interface{}{mock, "el nombre es: tinchogob"}
			},
		},
		{
			name: "Task/Do/return zero value ptr",
			tasks: func() []*task {
				return []*task{
					Task("p_1").Do(func() *task {
						return nil
					}),
					Task("p_2").After("p_1").Do(func(s1 *task) string {
						return "hola"
					}),
				}
			},
			results: func() []interface{} {
				return []interface{}{(*task)(nil), "hola"}
			},
		},
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

				if len(res) != len(results) {
					t.Fatalf("%s: expected %s results but got %s", c.name, len(results), len(res))
				}

				for i, r := range res {
					if r != results[i] {
						t.Errorf("%s: expected %v but got %v", c.name, results[i], r)
					}
				}
			}
		})
	}
}

func TestFlow(t *testing.T) {
	cases := []struct {
		name  string
		tasks func() []*task
		err   string
	}{
		{
			name: "Flow/After/1 dependency",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() time.Time {
						return time.Now()
					}),
					Task("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/2 dependency",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() time.Time {
						return time.Now()
					}),
					Task("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Task("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/3 dependency",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() time.Time {
						return time.Now()
					}),
					Task("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Task("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
					Task("four").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/4 dependency",
			tasks: func() []*task {
				return []*task{
					Task("one").Do(func() time.Time {
						return time.Now()
					}),
					Task("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Task("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
					Task("four").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Task("five").After("three", "four").Do(func(threeResult, fourResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tasks := c.tasks()
			res, err := Run(tasks...)
			if err != nil {
				if c.err == err.Error() {
					t.Fatalf("%s: expected %s error but got %s", c.name, c.err, err.Error())
				}
			} else {
				if c.err != "" {
					t.Fatalf("%s: expected %s ", c.name, c.err)
				}

				for i, task := range tasks {
					for _, dep := range task.dependencies {
						for j, depT := range tasks {
							if dep == depT.name && res[j].(time.Time).After(res[i].(time.Time)) {
								t.Errorf("%s: %s expected to be after %s ((%v is not after %v))", c.name, task.name, dep, res[j], res[i])
							}
						}
					}
				}
			}
		})
	}
}
