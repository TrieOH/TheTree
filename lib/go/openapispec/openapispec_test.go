package openapispec

import "testing"

func TestGeneratedID(t *testing.T) {
	cases := []struct{ spec, want string }{
		{"listActors", "ListActors"},
		{"getSetup", "GetSetup"},
		{"", ""},
		{"x", "X"},
	}
	for _, c := range cases {
		if got := GeneratedID(c.spec); got != c.want {
			t.Fatalf("GeneratedID(%q) = %q, want %q", c.spec, got, c.want)
		}
	}
}
