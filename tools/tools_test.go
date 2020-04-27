package tools

import (
	"testing"
)

func TestStr(t *testing.T) {
	var (
		in       = 43
		expected = "43"
	)
	out := Str(in)
	if out != expected {
		t.Errorf("Str(%d) = %s; expected %s", in, out, expected)
	}
}

func TestStr2(t *testing.T) {
	var (
		in       = 0.2143
		expected = "0.214300"
	)
	out := Str(in)
	if out != expected {
		t.Errorf("Str(%f) = %s; expected %s", in, out, expected)
	}
}

func TestInt(t *testing.T) {
	var (
		in       = "325"
		expected = 325
	)
	out := Int(in)
	if out != expected {
		t.Errorf("Int(%s) = %d; expected %d", in, out, expected)
	}
}

func TestLoadJSON(t *testing.T) {
	var (
		in = "../credential.json"
	)
	byteValue := LoadJSON(in)
	if byteValue[0] != 123 {
		t.Errorf("Error load json, expected first char is `{`, got %c ", byteValue[0])
	}
}
