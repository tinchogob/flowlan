package flowlan

import "fmt"

func ExampleRun() {

	res, _ := Run(Task("one").Do(func() string {
		return "one"
	}), Task("two").Do(func() string {
		return "two"
	}), Task("trhee").After("one","two").Do(func(one, two string) string {
		return one+"-"+two+"-"+"three"
	}))

	fmt.Println(res)
	//Output: [one two one-two-three]
}