package hasher

import "testing"

var HashTest = struct {
	input uint64
	want  string
	err   string
}{
	input: 18446744073709551615,
	want:  "pIrkgbKrQ8v",
}

func TestGenHash(t *testing.T) {
	res := GenHash(HashTest.input)
	if string(res) != HashTest.want {
		t.Errorf("%v: got - %v, expected - %v", HashTest.input, string(res), HashTest.want)
	}
}

func TestGenClear(t *testing.T) {
	res, _ := GenClear(HashTest.want)
	if res != HashTest.input {
		t.Errorf("%v: got - %v, expected - %v", HashTest.want, res, HashTest.input)
	}
}
