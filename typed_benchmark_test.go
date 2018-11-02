package flowlan

import "testing"

var res interface{}

func BenchmarkTyped(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Run(Step("one").Do(func() (string, error) {
			return "one", nil
		}), Step("two").Do(func() (string, string, error) {
			return "two", "two", nil
		}), Step("three").After("one", "two").Do(func(one, two, twoo string) (string, error) {
			return one + "-" + two + "-" + twoo + "-" + "three", nil
		}))
	}
}

func BenchmarkInterface(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Run(Step("one").IDo(func(map[string]interface{}) (interface{}, error) {
			return "one", nil
		}), Step("two").IDo(func(map[string]interface{}) (interface{}, error) {
			return "two", nil
		}), Step("three").After("one", "two").IDo(func(args map[string]interface{}) (interface{}, error) {
			one := args["one"].(string)
			two := args["two"].(string)
			return one + "-" + two + "-" + "three", nil
		}))
	}
}
