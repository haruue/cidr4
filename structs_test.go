package cidr4

import "testing"

func TestParseCIDRv4(t *testing.T) {
	tests := []string{
		"0.0.0.0/0",
		"192.0.2.0/24",
		"255.255.255.255/32",
	}
	for _, test := range tests {
		c, err := ParseCIDRv4(test)
		if err != nil {
			t.Errorf("expect: %s, err: %s", test, err.Error())
			continue
		}
		if test != c.String() {
			t.Errorf("expect: %s, actual: %s", test, c.String())
			continue
		}
	}
}

func TestCIDRv4_Standardize(t *testing.T) {
	type testCase struct {
		in  string
		out string
	}
	tests := []testCase{
		{in: "192.0.2.2/25", out: "192.0.2.0/25"},
		{in: "255.255.255.0/0", out: "0.0.0.0/0"},
	}
	for _, test := range tests {
		c, err := ParseCIDRv4(test.in)
		if err != nil {
			t.Errorf("in: %s, expect: %s, err: %s", test.in, test.out, err.Error())
			continue
		}
		if test.out != c.String() {
			t.Errorf("in: %s, expect: %s, actual: %s", test.in, test.out, c.String())
			continue
		}
	}
}

func TestCIDRv4_Contains(t *testing.T) {
	type testCase struct {
		container string
		subnet    string
		result    bool
	}
	tests := []testCase{
		{container: "192.0.2.0/24", subnet: "192.0.2.0/25", result: true},
		{container: "192.0.2.0/25", subnet: "192.0.2.0/25", result: true},
		{container: "192.0.2.0/24", subnet: "192.0.2.10/30", result: true},
		{container: "192.0.2.0/25", subnet: "192.0.2.140/30", result: false},
		{container: "192.0.2.0/25", subnet: "192.0.2.0/24", result: false},
		{container: "0.0.0.0/0", subnet: "255.255.255.255/32", result: true},
	}
	for _, test := range tests {
		c, err := ParseCIDRv4(test.container)
		if err != nil {
			t.Fatalf("failed to parse cidr: %s, err: %s", test.container, err.Error())
			return
		}
		s, err := ParseCIDRv4(test.subnet)
		if err != nil {
			t.Fatalf("failed to parse cidr: %s, err: %s", test.subnet, err.Error())
			return
		}
		result := c.Contains(s)
		if result != test.result {
			t.Errorf("(%s).Contains(%s): expect: %t, actual: %t", c.String(), s.String(), test.result, result)
			continue
		}
	}
}
