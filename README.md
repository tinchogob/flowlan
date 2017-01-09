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
	res, _ := flowlan.Run(flowlan.Task("times").Do(func() int {
		return 5
	}), flowlan.Task("message").Do(func() string, error {
		return "golang", nil
	}), flowlan.Task("repeat").After("times","message").Do(func(times int, msg string, msgErr error) string {
		var repeateadMsg string
        for i := 0; i < times; i++ {
            repeateadMsg += msg
        }
        return repeateadMsg
	}))

	fmt.Println(res)
	//Prints [5 golang <nil> golanggolanggolanggolanggolang]
	//times task returned- > 5
	//message task returned- > "golang", nil
	//repeat task returned golanggolanggolanggolanggolang
}
```

Built on top of go channels (with some reflection magic) as an excersise for go proverb:

_Do not communicate by sharing memory; instead, share memory by communicating._

TODO list
- General: flow validation
- General: error Handling
- General: cancellable contexts
- Do: support for variadic functions
- Do: support for timeout


