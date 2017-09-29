package flowlan

import "fmt"

func ExampleRun() {

	res, _ := Run(Task("one").Do(func() (string, error) {
		return "one", nil
	}), Task("two").Do(func() (string, string, error) {
		return "two", "two", nil
	}), Task("three").After("one", "two").Do(func(one, two, twoo string) (string, error) {
		return one + "-" + two + "-" + twoo + "-" + "three", nil
	}))

	fmt.Printf("one:%v two:%v three:%v ", res["one"], res["two"], res["three"])
	//Output: one:[one] two:[two two] three:[one-two-two-three]
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

	fmt.Printf("one:%v two:%v three:%v ", res["one"], res["two"], res["three"])
	//Output: one:one two:two three:one-two-three
}
