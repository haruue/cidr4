package cidrreverse

import "fmt"

func ensureCIDRv4(s string) CIDRv4 {
	c, err := ParseCIDRv4(s)
	if err != nil {
		panic(err)
	}
	return c
}

func ExampleCIDRv4Set_Minus() {
	cs := CIDRv4Set{}

	cs.Plus(ensureCIDRv4("0.0.0.0/0"))
	cs.Minus(ensureCIDRv4("192.0.2.0/24"))

	for _, c := range cs {
		fmt.Println(c.String())
	}

	// Output:
	// 0.0.0.0/1
	// 128.0.0.0/2
	// 192.0.0.0/23
	// 192.0.3.0/24
	// 192.0.4.0/22
	// 192.0.8.0/21
	// 192.0.16.0/20
	// 192.0.32.0/19
	// 192.0.64.0/18
	// 192.0.128.0/17
	// 192.1.0.0/16
	// 192.2.0.0/15
	// 192.4.0.0/14
	// 192.8.0.0/13
	// 192.16.0.0/12
	// 192.32.0.0/11
	// 192.64.0.0/10
	// 192.128.0.0/9
	// 193.0.0.0/8
	// 194.0.0.0/7
	// 196.0.0.0/6
	// 200.0.0.0/5
	// 208.0.0.0/4
	// 224.0.0.0/3
}
