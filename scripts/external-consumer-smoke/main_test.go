package main

import "testing"

func TestContainsReplaceDirective(t *testing.T) {
	for _, test := range []struct {
		name    string
		goMod   string
		replace bool
	}{
		{name: "published dependency", goMod: "module example.com/consumer\nrequire example.com/engine v0.1.0\n"},
		{name: "single replace", goMod: "replace example.com/engine => ../engine\n", replace: true},
		{name: "replace block", goMod: "replace (\nexample.com/engine => ../engine\n)\n", replace: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := containsReplaceDirective(test.goMod); got != test.replace {
				t.Fatalf("containsReplaceDirective() = %v, want %v", got, test.replace)
			}
		})
	}
}
