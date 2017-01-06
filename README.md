# Flowlan
#### simple flow orchestration for go

Flowlan provides a simple abstraction to execute an arbitrarily complex graph of dependable functions.
Simply declare a Task, its dependencies and its execution.

```golang
import (
    "fmt"
    
    "github.com/tinchogob/flowlan"
)

func main() {
	res, _ := flowlan.Run(flowlan.Task("one").Do(func() string {
		return "one"
	}), flowlan.Task("two").Do(func() string {
		return "two"
	}), flowlan.Task("trhee").After("one","two").Do(func(one, two string) string {
		return one+"-"+two+"-"+"three"
	}))

	fmt.Println(res)
	//Prints [one two one-two-three]
}
```

Built on top of go channels (with some reflection magic) as an excersise for go proverb:

_Do not communicate by sharing memory; instead, share memory by communicating._

Caution: calling a function via reflection in go is [slow!](https://github.com/golang/go/issues/7818) and this makes heavy use of it

TODO list
- Do: support for variadic functions
- Do: support for timeout
- General: error Handling
- General: cancellable contexts

