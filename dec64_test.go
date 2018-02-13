package dec64

import (
	"math"
	"strconv"
	"testing"
)

func testOneDec(t *testing.T, s, refs string, ref int64) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%s Result is %d should be %d", s, d, ref)
	}
	if refs != d.String() {
		t.Errorf("String is %s should be %s", d.String(), refs)
	}
	f, _ := strconv.ParseFloat(s, 64)
	if math.Abs(f-Float64(d)) > 0.000000000000001 {
		t.Errorf("Float64 is %g should be %g", Float64(d), f)
	}
}

func TestParse(t *testing.T) {
	testOneDec(t, "1", "1", 256)
	testOneDec(t, "-1", "-1", -256)
	testOneDec(t, "345", "345", 345*256)
	testOneDec(t, ".6789", ".6789", 1738236)
	testOneDec(t, "2.05", "2.05", 52734)
	testOneDec(t, "300201", "300201", 76851456)
	testOneDec(t, "100", "100", 258)
	testOneDec(t, ".09", ".09", 9*256+254)
	testOneDec(t, ".007", ".007", 7*256+256-3)
	testOneDec(t, "0", "0", 0)
	testOneDec(t, "1E-8", ".00000001", 1*256+256-8)
	testOneDec(t, "3E-8", ".00000003", 3*256+256-8)
	testOneDec(t, "1.2E-7", ".00000012", 12*256+256-8)
}
