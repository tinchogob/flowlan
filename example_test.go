package flowlan

import "fmt"

func ExampleRun() {

	var r1, r2, r3 string
	Run(Step("one").Do(func() (string, error) {
		return "one", nil
	}), Step("two").Do(func() (string, string, error) {
		return "two", "two", nil
	}), Step("three").After("one", "two").Do(func(one, two, twoo string) (string, error) {
		r1 = one
		r2 = two
		r3 = one + "-" + two + "-" + twoo + "-" + "three"
		return "", nil
	}))

	fmt.Printf("one:%v two:%v three:%v ", r1, r2, r3)
	//Output: one:one two:two three:one-two-two-three
}

func ExampleIRun() {

	var one, two, three string
	Run(Step("one").IDo(func(map[string]interface{}) (interface{}, error) {
		return "one", nil
	}), Step("two").IDo(func(map[string]interface{}) (interface{}, error) {
		return "two", nil
	}), Step("three").After("one", "two").IDo(func(args map[string]interface{}) (interface{}, error) {
		one = args["one"].(string)
		two = args["two"].(string)
		three = one + "-" + two + "-" + "three"
		return nil, nil
	}))

	fmt.Printf("one:%v two:%v three:%v ", one, two, three)
	//Output: one:one two:two three:one-two-three
}
