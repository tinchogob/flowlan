package flowlan

import (
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Debug = true
	cases := []struct {
		name  string
		tasks func() []*Task
		err   string
	}{
		/*
			{
				name: "Step/Do/no action",
				tasks: func() []*Task {
					return []*Task{Step("Tincho")}
				},
			},
			{
				name: "Step/Do/no return values",
				tasks: func() []*Task {
					return []*Task{
						Step("one").Do(func() {
							time.Sleep(time.Microsecond)
						}),
						Step("two").After("one").Do(func() {
							time.Sleep(time.Microsecond)
						}),
					}
				},
			},
			{
				name: "Step/Do/return values",
				tasks: func() []*Task {
					return []*Task{
						Step("one").Do(func() int {
							return 1
						}),
						Step("two").After("one").Do(func(one int) int {
							if one != 1 {
								t.Fail()
							}
							return one + 2
						}),
					}
				},
			},
			{
				name: "Task/Do/multiple return values",
				tasks: func() []*Task {
					return []*Task{
						Step("multiple").Do(func() (int, int) {
							return 1, 2
						}),
						Step("last").After("multiple").Do(func(one, two int) {
							if one+two != 3 {
								t.Fail()
							}
						}),
					}
				},
			},*/
		{
			name: "Task/Do/error return value",
			tasks: func() []*Task {
				return []*Task{
					Step("first").Do(func() error {
						return errors.New("lala")
					}),
					Step("last").After("first").Do(func(e1 error) {
						if e1.Error() != "lala" {
							t.Fail()
						}
					}),
				}
			},
		},
		/*
			{
				name: "Task/Do/zero return value error",
				tasks: func() []*Task {
					return []*Task{
						Step("zre_1").Do(func() error {
							return nil
						}),
						Step("zre_2").After("zre_1").Do(func(e1 error) int {
							return 1
						}),
					}
				},

			},
			{
				name: "Task/Do/return value ptr",
				tasks: func() []*Task {
					return []*Task{
						Step("s_1").Do(func() *string {
							str := "tinchogob"
							return &str
						}),
						Step("s_2").After("s_1").Do(func(s1 *Task) string {
							return "el nombre es: " + s1.Name
						}),
					}
				},

			},
			{
				name: "Task/Do/return zero value ptr",
				tasks: func() []*Task {
					return []*Task{
						Step("p_1").Do(func() *Task {
							return nil
						}),
						Step("p_2").After("p_1").Do(func(s1 *Task) string {
							return "hola"
						}),
					}
				},

			},
			{
				name: "Task/Do/return errors",
				tasks: func() []*Task {
					return []*Task{
						Step("p_1").Do(func() (*Task, error) {
							return nil, io.ErrUnexpectedEOF
						}),
						Step("p_2").After("p_1").Do(func(s1 *Task, e1 error) int {
							return 5
						}),
					}
				},

			},

			{
				name: "Flow/After/wrong dependency",
				tasks: func() []*Task {
					return []*Task{
						Step("one").Do(func() int {
							return 1
						}),
						Step("two").After("chabon").Do(func(oneResult int) int {
							return 1 + oneResult
						}),
					}
				},
				err: "error",
			},

			{
				name: "Flow/After/two many dependencies",
				tasks: func() []*Task {
					return []*Task{
						Step("one").Do(func() int {
							return 1
						}),
						Step("two").After("chabon").Do(func() int {
							return 1
						}),
					}
				},
				err: "error",
			},
			{
				name: "Flow/After/self dependency",
				tasks: func() []*Task {
					return []*Task{
						Step("one").Do(func() int {
							return 1
						}),
						Step("two").After("two").Do(func() int {
							return 1
						}),
					}
				},
				err: "error",
			},
			{
				name: "Flow/After/circular dependencies",
				tasks: func() []*Task {
					return []*Task{
						Step("one").After("three").Do(func() int {
							return 1
						}),
						Step("two").After("one").Do(func() int {
							return 1
						}),
						Step("three").After("two").Do(func() int {
							return 1
						}),
					}
				},
				err: "error",
			},
			{
				name: "Flow/After/deadlock dependencies",
				tasks: func() []*Task {
					return []*Task{
						Step("one").After("two").Do(func() int {
							return 1
						}),
						Step("two").After("one").Do(func() int {
							return 1
						}),
					}
				},
				err: "error",
			},
		*/
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Run(c.tasks()...)
			if err != nil {
				if c.err == err.Error() {
					t.Fatalf("%s: expected %s error but got %s", c.name, c.err, err.Error())
				}
			} else {
				if c.err != "" {
					t.Fatalf("%s: expected %s ", c.name, c.err)
				}
			}
		})
	}
}

func TestFlow(t *testing.T) {
	cases := []struct {
		name  string
		tasks func() []*Task
		err   string
	}{
		{
			name: "Flow/After/1 dependency",
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() time.Time {
						return time.Now()
					}),
					Step("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/2 dependency",
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() time.Time {
						return time.Now()
					}),
					Step("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Step("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/3 dependency",
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() time.Time {
						return time.Now()
					}),
					Step("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Step("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
					Step("four").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
		{
			name: "Flow/After/4 dependency",
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() time.Time {
						return time.Now()
					}),
					Step("two").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Step("three").After("one", "two").Do(func(oneResult, twoResult time.Time) time.Time {
						return time.Now()
					}),
					Step("four").After("one").Do(func(oneResult time.Time) time.Time {
						return time.Now()
					}),
					Step("five").After("three", "four").Do(func(threeResult, fourResult time.Time) time.Time {
						return time.Now()
					}),
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tasks := c.tasks()
			err := Run(tasks...)
			if err != nil {
				if c.err == err.Error() {
					t.Fatalf("%s: expected %s error but got %s", c.name, c.err, err.Error())
				}
			} else {
				if c.err != "" {
					t.Fatalf("%s: expected %s ", c.name, c.err)
				}

				// for _, task := range tasks {
				// 	for _, dep := range task.in {
				// 		for _, depT := range tasks {
				// 			if dep.name == depT.name && res[depT.name].(time.Time).After(res[task.name].(time.Time)) {
				// 				t.Errorf("%s: %s expected to be after %s ((%v is not after %v))", c.name, task.name, dep.name, res[depT.name], res[task.name])
				// 			}
				// 		}
				// 	}
				// }
			}
		})
	}
}
