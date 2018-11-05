package flowlan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Debug = false
	cases := []struct {
		name  string
		tasks func() []*Task
		ctx   func() context.Context
		err   string
	}{
		{
			name: "Step/Do/no action",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{Step("Tincho")}
			},
		},
		{
			name: "Step/Do/no return values",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() {}),
					Step("two").After("one").Do(func() {}),
				}
			},
		},
		{
			name: "Step/Do/return values",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() int {
						return 1
					}),
					Step("two").After("one").Do(func(one int) {
						if one != 1 {
							t.Fail()
						}
					}),
				}
			},
		},
		{
			name: "Task/Do/multiple return values",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("multiple").Do(func() (int, int) {
						return 1, 2
					}),
					Step("last").After("multiple").Do(func(one, two int) {
						if one != 1 {
							t.Fail()
						}
						if two != 2 {
							t.Fail()
						}
					}),
				}
			},
		},
		{
			name: "Step/Do/return interface",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() io.Reader {
						return strings.NewReader("Hello, Reader!")
					}),
					Step("two").After("one").Do(func(in io.Reader) {
						b := make([]byte, 8)
						for {
							_, err := in.Read(b)
							if err == io.EOF {
								return
							}
						}
					}),
				}
			},
		},
		{
			name: "Step/Do/return nil interface",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() io.Reader {
						return nil
					}),
					Step("two").After("one").Do(func(in io.Reader) {
						if in != nil {
							t.Fail()
						}
					}),
				}
			},
		},
		{
			name: "Task/Do/error return value",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("first").Do(func() error {
						return errors.New("lala")
					}),
					Step("last").After("first").Do(func(e1 error) {
						if e1 == nil {
							fmt.Println("no deberia ser nil")
							t.Fail()
						}
						if e1.Error() != "lala" {
							t.Fail()
						}
					}),
				}
			},
		},
		{
			name: "Task/Do/return nil error",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("zre_1").Do(func() error {
						return nil
					}),
					Step("zre_2").After("zre_1").Do(func(e1 error) {
						if e1 != nil {
							fmt.Println("deberia ser nil")
							t.Fail()
						}
					}),
				}
			},
		},
		{
			name: "Task/Do/return value ptr",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("s_1").Do(func() *string {
						str := "tinchogob"
						return &str
					}),
					Step("s_2").After("s_1").Do(func(s1 *string) string {
						return "el nombre es: " + *s1
					}),
				}
			},
		},
		{
			name: "Task/Do/return zero value ptr",
			ctx: func() context.Context {
				return context.Background()
			},
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
			ctx: func() context.Context {
				return context.Background()
			},
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
			name: "Flow/Do/many args",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() (string, error) {
						return "one", nil
					}), Step("two").Do(func() (string, string, error) {
						return "two", "two", nil
					}), Step("three").After("one", "two").Do(func(one string, oneErr error, two, twoo string, toErr error) (string, error) {
						return one + "-" + two + "-" + twoo + "-" + "three", nil
					}),
				}
			},
		},
		{
			name: "Flow/Do/with context",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() (string, error) {
						return "one", nil
					}), Step("two").Do(func() (string, string, error) {
						return "two", "two", nil
					}), Step("three").After("one", "two").Do(func(one string, oneErr error, two, twoo string, toErr error) (string, error) {
						return one + "-" + two + "-" + twoo + "-" + "three", nil
					}),
				}
			},
		},
		{
			name: "Flow/After/wrong dependency",
			ctx: func() context.Context {
				return context.Background()
			},
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
			err: "invalid task definition: two unknown dependency chabon",
		},
		{
			name: "Flow/After/two many wrong dependencies",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() int {
						return 1
					}),
					Step("two").After("chabon", "chabona").Do(func() int {
						return 1
					}),
				}
			},
			err: "invalid task definition: two unknown dependency chabon",
		},
		{
			name: "Flow/After/self dependency",
			ctx: func() context.Context {
				return context.Background()
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() int {
						return 1
					}),
					Step("two").After("two").Do(func(a int) int {
						return 1
					}),
				}
			},
			err: "invalid task definition: circular dependency on two",
		},
		{
			name: "Flow/Context/done context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() int {
						return 1
					}),
					Step("two").After("two").Do(func(a int) int {
						return 1
					}),
				}
			},
			err: "context canceled",
		},
		{
			name: "Flow/Context/timeouts",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*5)
				return ctx
			},
			tasks: func() []*Task {
				return []*Task{
					Step("one").Do(func() int {
						<-time.After(time.Millisecond * 10)
						return 2
					}),
					Step("two").After("one").Do(func(a int) int {
						return 1
					}),
				}
			},
			err: "context deadline exceeded",
		},
	}

	for _, c := range cases {
		log(c.name)
		t.Run(c.name, func(t *testing.T) {
			err := Run(c.ctx(), c.tasks()...)
			if err != nil {
				if c.err != err.Error() {
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
