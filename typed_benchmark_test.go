package flowlan

import "testing"

var res interface{}

func BenchmarkTyped(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r, _ := Run(Task("one").Do(func() (string, error) {
			return "one", nil
		}), Task("two").Do(func() (string, string, error) {
			return "two", "two", nil
		}), Task("three").After("one", "two").Do(func(one, two, twoo string) (string, error) {
			return one + "-" + two + "-" + twoo + "-" + "three", nil
		}))

		res = r
	}
}

func BenchmarkInterface(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r, _ := Run(Task("one").IDo(func(map[string]interface{}) (interface{}, error) {
			return "one", nil
		}), Task("two").IDo(func(map[string]interface{}) (interface{}, error) {
			return "two", nil
		}), Task("three").After("one", "two").IDo(func(args map[string]interface{}) (interface{}, error) {
			one := args["one"].(string)
			two := args["two"].(string)
			return one + "-" + two + "-" + "three", nil
		}))

		res = r
	}
}
