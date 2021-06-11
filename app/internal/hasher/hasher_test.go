package hasher

import "testing"

var HashTest = struct {
	num uint64
	str string
	err string
}{
	num: 18446744073709551615,
	str: "pIrkgbKrQ8v",
}

func TestGenHash(t *testing.T) {
	res := GenHash(HashTest.num)
	if string(res) != HashTest.str {
		t.Errorf("%v: got - %v, expected - %v", HashTest.num, string(res), HashTest.str)
	}
}

func TestGenClear(t *testing.T) {
	res, _ := GenClear(HashTest.str)
	if res != HashTest.num {
		t.Errorf("%v: got - %v, expected - %v", HashTest.str, res, HashTest.num)
	}
}
