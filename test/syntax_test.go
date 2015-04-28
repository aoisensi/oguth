package test

import (
	"testing"

	"github.com/aoisensi/oguth"
)

func TestGeneratedCodeSyntax(t *testing.T) {
	//https://tools.ietf.org/html/rfc6749#appendix-A.11

	for i := 0; i < 100; i++ {
		code := oguth.DefaultAuthCodeGenerator()
		if !VSCHARSyntax(code) {
			t.Fail()
		}
	}
}

func VSCHARSyntax(s string) bool {
	for _, c := range s {
		if c < 0x20 || c > 0x7e {
			return false
		}
	}
	return true
}
