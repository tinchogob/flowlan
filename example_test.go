package flowlan

import (
	"fmt"
)

func ExampleFlowlan() {

	t1 := Task("primera").Do(func() (string, error) {
		return "el res de primera", nil
	})

	t2 := Task("segunda").Do(func() (string, error) {
		return "", nil
	})

	t3 := Task("tercera").After("segunda").Do(func(r1 string) (string, error) {
		return "", nil
	})

	t4 := Task("cuarta").After("segunda").Do(func(s string) (string, error) {
		fmt.Printf("exec cuarta con len de args: %d\n", len(s))
		return "%"+s+"modificado", nil
	})

	t5 := Task("quinta").After("caca").Do(func(s string) (string, error) {
		return "el res de quinta", nil
	})

	res, _ := Run(t1, t2, t3, t4, t5)

	fmt.Println(res)
	//Output: hola

}
