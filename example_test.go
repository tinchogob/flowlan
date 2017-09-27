package flowlan

import "fmt"

func ExampleRun() {

	res, _ := Run(Task("one").Do(func() (string, error) {
		return "one", nil
	}), Task("two").Do(func() (string, error) {
		return "two", nil
	}), Task("trhee").After("one", "two").Do(func(one, two string) (string, error) {
		return one + "-" + two + "-" + "three", nil
	}))

	fmt.Println(res)
	//Output: [one two one-two-three]
}

func ExampleIRun() {

	res, _ := Run(Task("one").IDo(func(map[string]interface{}) (interface{}, error) {
		return "one", nil
	}), Task("two").IDo(func(map[string]interface{}) (interface{}, error) {
		return "two", nil
	}), Task("three").After("one", "two").IDo(func(args map[string]interface{}) (interface{}, error) {
		one := args["one"].(string)
		two := args["two"].(string)
		return one + "-" + two + "-" + "three", nil
	}))

	fmt.Println(res)
	//Output: map[one:one two:two three:one-two-three]
}
